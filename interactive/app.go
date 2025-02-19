package main

import (
	"github.com/chenmuyao/go-bootcamp/internal/events"
	"github.com/chenmuyao/go-bootcamp/pkg/grpcx"
)

type App struct {
	consumers []events.Consumer
	server    *grpcx.Server
}
