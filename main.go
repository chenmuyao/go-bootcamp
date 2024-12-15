package main

import (
	"github.com/chenmuyao/go-bootcamp/config"
	"go.uber.org/zap"
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
	initLogger()

	server := InitWebServer()

	server.Run(":8081")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Debug("hello world")
}
