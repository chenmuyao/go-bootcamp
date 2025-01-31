package main

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/config"
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
	intrRepo.BatchSetTopLike(ctx, "article", 1000)

	app.server.Run(":8081")
}
