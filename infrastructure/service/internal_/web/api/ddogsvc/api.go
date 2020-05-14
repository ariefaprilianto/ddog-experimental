package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ariefaprilianto/ddog-experimental/infrastructure/config"
	metric "github.com/ariefaprilianto/ddog-experimental/infrastructure/metric/definitions"
	myrouter "github.com/ariefaprilianto/ddog-experimental/infrastructure/service/internal_/router"
	"github.com/ariefaprilianto/ddog-experimental/lib/common/response"
	"github.com/julienschmidt/httprouter"
)

const (
	HTTPGenericSuccess          = 200
	HTTPCodeBadRequest          = 400
	HTTPForbiddenResource       = 403
	HTTPCodeInternalServerError = 500
)

type Metric struct {
	DDogSvcMetric metric.MetricInterface
}

// API is the api struct
type API struct {
	Cfg    *config.MainConfig
	Metric *Metric
}

type controlledBehaviour struct {
	Err             error
	LatencyInSecond int
}

// New is the api initializer
func New(this *API) *API {
	return &API{
		Cfg:    this.Cfg,
		Metric: this.Metric,
	}
}

// Register will register the api structure
func (a *API) Register() {
	router := myrouter.New(&myrouter.Options{Timeout: a.Cfg.API.DefaultTimeout, Prefix: a.Cfg.API.NormalPrefix})
	router.GET("/accounts", a.Accounts)
	router.GET("/customers", a.Customers)
}

// Accounts handle accounts endpoint
func (a *API) Accounts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *response.JSONResponse {
	behaviour, err := parseControlledBehaviour(r)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(fmt.Sprintf("%s error - %s", "Accounts", err.Error()))
	}

	if behaviour.LatencyInSecond > 0 {
		timeDuration := time.Duration(behaviour.LatencyInSecond)
		time.Sleep(timeDuration * time.Second)

		log.Println("Latency: ", behaviour.LatencyInSecond)
	}
	if behaviour.Err != nil {
		return response.NewJSONResponse().SetError(behaviour.Err).SetMessage(fmt.Sprintf("%s error - %s", "Accounts", behaviour.Err.Error()))
	}

	timeDuration := time.Duration(behaviour.LatencyInSecond)
	if behaviour.LatencyInSecond > 0 {
		time.Sleep(timeDuration * time.Second)
	}

	return response.NewJSONResponse().SetData("Succeeded")
}

// Customers handle customers endpoint
func (a *API) Customers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *response.JSONResponse {
	behaviour, err := parseControlledBehaviour(r)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage(fmt.Sprintf("%s error - %s", "Customers", response.ErrInternalServerError.Error()))
	}

	if behaviour.LatencyInSecond > 0 {
		timeDuration := time.Duration(behaviour.LatencyInSecond)
		time.Sleep(timeDuration * time.Second)

		log.Println("Latency: ", behaviour.LatencyInSecond)
	}
	if behaviour.Err != nil {
		return response.NewJSONResponse().SetError(behaviour.Err).SetMessage(fmt.Sprintf("%s error - %s", "Customers", behaviour.Err.Error()))
	}

	return response.NewJSONResponse().SetData("Succeeded")
}

func parseControlledBehaviour(r *http.Request) (b controlledBehaviour, err error) {
	b = controlledBehaviour{}

	rawStatusCode := r.URL.Query().Get("status_code")
	if len(rawStatusCode) > 0 {
		statusCode, e := strconv.Atoi(rawStatusCode)
		if e != nil {
			err = e
			return
		}

		switch statusCode {
		case HTTPCodeBadRequest:
			b.Err = response.ErrBadRequest
			break
		case HTTPForbiddenResource:
			b.Err = response.ErrForbiddenResource
			break
		case HTTPGenericSuccess:
			b.Err = nil
			break
		default:
			b.Err = response.ErrInternalServerError
			break
		}
	}

	rawLatencyInSecond := r.URL.Query().Get("latency")
	if len(rawLatencyInSecond) > 0 {
		latencyInSecond, e := strconv.Atoi(rawLatencyInSecond)
		if e != nil {
			err = e
			return
		}
		b.LatencyInSecond = latencyInSecond
	}

	return
}
