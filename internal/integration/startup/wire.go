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

var thirdPartySet = wire.NewSet(
	InitRedis,
	InitDB,
	InitLogger,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third-party dependencies
		thirdPartySet,

		// DAO
		dao.NewUserDAO, dao.NewAsyncSMSDAO, dao.NewArticleDAO,

		// Cache
		rediscache.NewCodeRedisCache,
		rediscache.NewUserRedisCache,
		rediscache.NewArticleRedisCache,
		// ioc.InitCodeLocalCache,
		// ioc.InitUserLocalCache,

		// Repo
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewAsyncSMSRepository,
		repository.NewArticleRepository,

		// Services
		ioc.InitSMSService,
		service.NewCodeService,
		service.NewUserService,
		service.NewArticleService,

		// handler
		web.NewUserHandler,
		NewDummyGiteaHandler,
		ijwt.NewRedisJWTHandler,
		web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return &gin.Engine{}
}

func InitArticleHandler(dao dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		rediscache.NewArticleRedisCache,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
