package main

import (
	"context"
	"net/http"
	"time"

	"github.com/chenmuyao/go-bootcamp/interactive/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}

func main() {
	config.InitConfig("config/dev.yaml")

	app := InitApp()

	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	intrRepo := InitInteractiveRepo()
	err := intrRepo.BatchSetTopLike(ctx, "article", 1000)
	if err != nil {
		panic(err)
	}
	err = app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
}
