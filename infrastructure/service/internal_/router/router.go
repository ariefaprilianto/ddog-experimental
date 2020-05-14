package router

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	metric "github.com/ariefaprilianto/ddog-experimental/infrastructure/metric/definitions"
	"github.com/ariefaprilianto/ddog-experimental/lib/common/response"
	"github.com/felixge/httpsnoop"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func init() {
	HttpRouter = httprouter.New()
}

type MyRouter struct {
	Httprouter     *httprouter.Router
	WrappedHandler http.Handler
	Options        *Options
}

type Options struct {
	Prefix  string
	Timeout int
}

type WrittenResponseWriter struct {
	http.ResponseWriter
	written bool
}

func (w *WrittenResponseWriter) WriteHeader(status int) {
	w.written = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *WrittenResponseWriter) Write(b []byte) (int, error) {
	w.written = true
	return w.ResponseWriter.Write(b)
}

func (w *WrittenResponseWriter) Written() bool {
	return w.written
}

var HttpRouter *httprouter.Router

// WrapperHandler used to wrap web handler
func WrapperHandler(metric metric.MetricInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writtenResponseWriter := &WrittenResponseWriter{
			ResponseWriter: w,
			written:        false,
		}
		w = writtenResponseWriter

		// metric data
		start := time.Now()

		// CaptureMetrics wraps the given handler, executes it with the given w and r, and
		// returns the metrics captured from it within processing time from start to finish.
		m := httpsnoop.CaptureMetrics(HttpRouter, w, r)

		// capture route-path from request header for metric tagging
		urlPathTag := r.Header.Get("routePath")
		if urlPathTag == "" {
			urlPathTag = "UNKNOWN"
		}

		// define datadog metric tags
		tags := []string{
			"via:http",
			fmt.Sprintf("url_path:%s", urlPathTag),
			fmt.Sprintf("url:%s", r.URL.Path),
			fmt.Sprintf("resp_code:%d", m.Code),
		}

		// metric submission to ddog agent, in here comm. will be UDP, so we no need to wait the response back
		go metric.Histogram("http_router", start, tags)
	})
}

func New(o *Options) *MyRouter {
	myrouter := &MyRouter{
		Options:    o,
		Httprouter: HttpRouter,
	}
	return myrouter
}

type Handle func(http.ResponseWriter, *http.Request, httprouter.Params) *response.JSONResponse

func (mr *MyRouter) GET(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.GET(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) GETFile(path string, handle httprouter.Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.GET(fullPath, handle)
}

func (mr *MyRouter) POST(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.POST(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) PUT(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.PUT(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) PATCH(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.PATCH(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) DELETE(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.DELETE(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) OPTIONS(path string, handle Handle) {
	fullPath := mr.Options.Prefix + path
	log.Println(fullPath)
	mr.Httprouter.OPTIONS(fullPath, mr.handleNow(fullPath, handle))
}

func (mr *MyRouter) ServeFiles(path string, root http.FileSystem) {
	mr.Httprouter.ServeFiles(path, root)
}

func (mr *MyRouter) TestHack(fullPath string, handle Handle) httprouter.Handle {
	return mr.handleNow(fullPath, handle)
}

func (mr *MyRouter) handleNow(fullPath string, handle Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		t := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(mr.Options.Timeout))

		defer cancel()

		ctx = context.WithValue(ctx, "HTTPParams", ps)

		r.Header.Set("routePath", fullPath)
		r = r.WithContext(ctx)

		respChan := make(chan *response.JSONResponse)
		go func() {
			defer panicRecover(r, fullPath)
			resp := handle(w, r, ps)
			respChan <- resp
		}()

		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
				_, err := w.Write([]byte("timeout")) //TODO: should be custom response
				if err != nil {
					log.Println(err)
				}
			}
		case resp := <-respChan:
			if resp != nil {
				resp.SetLatency(time.Since(t).Seconds() * 1000)
				log.WithFields(log.Fields{
					"Resp-Status-Code": resp.StatusCode,
					"Latency":          resp.Latency,
					"Request-URI":      r.URL.RequestURI(),
				}).Info("Request processed")
				resp.Send(w)
			} else {
				if w, ok := w.(*WrittenResponseWriter); ok && !w.Written() {
					log.Println("Error nil response from the handler")
					w.WriteHeader(http.StatusInternalServerError)
					_, err := w.Write([]byte(""))
					if err != nil {
						log.Println(err)
					}
				}
			}
		}

		return
	}
}

func GetHttpParam(ctx context.Context, name string) string {
	ps := ctx.Value("HTTPParams").(httprouter.Params)
	return ps.ByName(name)
}

func panicRecover(r *http.Request, path string) {
	if err := recover(); err != nil {
		stackTrace := string(debug.Stack())
		log.Println("got panic in api handler, [Path] %s, [err] %v, stacktrace", path, err, stackTrace)
	}
}
