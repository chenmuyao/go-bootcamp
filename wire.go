//go:build wireinject
// +build wireinject

package main

import (
	"github.com/chenmuyao/go-bootcamp/internal/events/article"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/google/wire"
)

var interactiveSvcSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	rediscache.NewInteractiveRedisCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService,
)

func InitWebServer() *App {
	wire.Build(
		// third-party dependencies
		ioc.InitRedis,
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitSaramaClient,
		ioc.InitSyncProducer,

		interactiveSvcSet,

		article.NewSaramaSyncProducer,
		article.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,

		// DAO
		dao.NewUserDAO,
		dao.NewAsyncSMSDAO,
		dao.NewArticleDAO,

		// Cache
		rediscache.NewCodeRedisCache,
		rediscache.NewUserRedisCache,
		rediscache.NewArticleRedisCache,
		// ioc.InitCodeLocalCache,
		// ioc.InitUserLocalCache,
		ioc.InitTopArticlesCache,

		// Repo
		repository.NewUserRepository,
		repository.NewCodeRepository,
		repository.NewAsyncSMSRepository,
		repository.NewArticleRepository,

		// Services
		ioc.InitSMSService,
		service.NewCodeService,
		service.NewUserService,
		ioc.InitGiteaService,
		service.NewArticleService,

		// handler
		web.NewUserHandler,
		web.NewOAuth2GiteaHandler,
		ijwt.NewRedisJWTHandler,
		web.NewArticleHandler,

		ioc.InitGinMiddlewares,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	return new(App)
}
