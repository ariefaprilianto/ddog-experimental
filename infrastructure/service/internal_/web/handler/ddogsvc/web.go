package handler

import (
	"log"

	"github.com/ariefaprilianto/ddog-experimental/infrastructure/config"
	myrouter "github.com/ariefaprilianto/ddog-experimental/infrastructure/service/internal_/router"
	api "github.com/ariefaprilianto/ddog-experimental/infrastructure/service/internal_/web/api/ddogsvc"
	"gopkg.in/tokopedia/grace.v1"
)

//Handler is Web requests handler struct
type Handler struct {
	Cfg         *config.MainConfig
	Metric      *api.Metric
	listenErrCh chan error
}

// New is the web handler initializer
func New(this *Handler) *Handler {
	a := &api.API{Cfg: this.Cfg, Metric: this.Metric}
	api.New(a).Register()
	return this
}

//Run is to run the web apis
func (h *Handler) Run() {
	log.Printf("Listening on %s", h.Cfg.Server.Port)
	h.listenErrCh <- grace.Serve(h.Cfg.Server.Port, myrouter.WrapperHandler(h.Metric.DDogSvcMetric))
}

//ListenError will lister the error
func (h *Handler) ListenError() <-chan error {
	return h.listenErrCh
}
