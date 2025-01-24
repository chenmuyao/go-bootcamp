package main

import (
	"github.com/chenmuyao/go-bootcamp/internal/events"
	"github.com/gin-gonic/gin"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
