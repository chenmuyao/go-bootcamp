//go:build wireinject
// +build wireinject

package main

import (
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third-party dependencies
		ioc.InitRedis, ioc.InitDB,

		// DAO
		dao.NewUserDAO,

		// Cache
		rediscache.NewCodeRedisCache, rediscache.NewUserRedisCache,

		// Repo
		repository.NewUserRepository, repository.NewCodeRepository,

		// Services
		ioc.InitSMSService,
		service.NewCodeService,
		service.NewUserService,

		// handler
		web.NewUserHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return &gin.Engine{}
}
