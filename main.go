package main

import (
	"context"
	"net/http"
	"time"

	"github.com/chenmuyao/go-bootcamp/config"
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

	initPrometheus()

	app := InitWebServer()

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

	app.server.Run(":8081")
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8082", nil)
	}()
}
