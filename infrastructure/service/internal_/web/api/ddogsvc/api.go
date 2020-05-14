package api

import (
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
	time.Sleep(5 * time.Second)

	behaviour, err := parseControlledBehaviour(r)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError)
	}
	if behaviour.Err != nil {
		return response.NewJSONResponse().SetError(behaviour.Err)
	}

	return response.NewJSONResponse().SetData("You've been logged out successfully")
}

// Customers handle customers endpoint
func (a *API) Customers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *response.JSONResponse {
	return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage("Customers call error occured")

	// return response.NewJSONResponse().SetData("You've been logged out successfully")
}

func parseControlledBehaviour(r *http.Request) (b controlledBehaviour, err error) {
	b = controlledBehaviour{}

	statusCode, err := strconv.Atoi(r.URL.Query().Get("status_code"))
	if err != nil {
		return
	}

	latencyInSecond, err := strconv.Atoi(r.URL.Query().Get("latency"))
	if err != nil {
		return
	}

	b.LatencyInSecond = latencyInSecond

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

	return
}
