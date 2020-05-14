package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ariefaprilianto/ddog-experimental/infrastructure/config"
	"github.com/ariefaprilianto/ddog-experimental/infrastructure/metric/implementation/datadog"
	api "github.com/ariefaprilianto/ddog-experimental/infrastructure/service/internal_/web/api/ddogsvc"
	handler "github.com/ariefaprilianto/ddog-experimental/infrastructure/service/internal_/web/handler/ddogsvc"
	"github.com/ariefaprilianto/ddog-experimental/lib/common/env"
	_ "github.com/tokopedia/dexter/profx/integration"
)

func main() {
	os.Exit(Main())
}

// SvcName where you should put your custom service name here to distinguish stored ddog metric
const SvcName = "ddogsvc"

// Main is the main function
func Main() int {
	// configuration init
	cfg := getConfig()

	log.Printf("%s started,\n cfg=%+v", cfg.Server.Name, cfg) //message will not appear unless run with -debug switch

	// metric initialization
	datadogClient := datadog.New(SvcName, env.Get(), cfg.Datadog.Endpoint)
	metric := &api.Metric{
		DDogSvcMetric: datadogClient,
	}

	// init server
	h := handler.Handler{Cfg: cfg, Metric: metric}
	server := handler.New(&h)
	fmt.Println(fmt.Printf("%+v", h))
	go server.Run()

	// catch terminal os signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-term:
		log.Println("Exiting gracefully...", s)
	case err := <-server.ListenError():
		log.Println("Error starting web server, exiting gracefully:", err)
	}

	return 0
}

//EarlyExit from the app
func earlyExit(flag bool) {
	if flag {
		os.Exit(0)
		return
	}
}

func getConfig() *config.MainConfig {
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, SvcName, "main")
	return cfg
}
