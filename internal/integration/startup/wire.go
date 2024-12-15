//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third-party dependencies
		InitRedis,
		ioc.InitDB,
		InitLogger,

		// DAO
		dao.NewUserDAO, dao.NewAsyncSMSDAO,

		// Cache
		rediscache.NewCodeRedisCache, rediscache.NewUserRedisCache,
		// ioc.InitCodeLocalCache, ioc.InitUserLocalCache,

		// Repo
		repository.NewUserRepository, repository.NewCodeRepository,
		repository.NewAsyncSMSRepository,

		// Services
		ioc.InitSMSService,
		service.NewCodeService,
		service.NewUserService,

		// handler
		web.NewUserHandler,
		NewDummyGiteaHandler,
		ijwt.NewRedisJWTHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return &gin.Engine{}
}
