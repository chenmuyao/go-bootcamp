package main

import (
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

	app.server.Run(":8081")
}
