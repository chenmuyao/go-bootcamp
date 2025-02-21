//go:build wireinject
// +build wireinject

package startup

import (
	intrRepository "github.com/chenmuyao/go-bootcamp/interactive/repository"
	intrRediscache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache/rediscache"
	intrDao "github.com/chenmuyao/go-bootcamp/interactive/repository/dao"
	intrService "github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/internal/events/article"
	"github.com/chenmuyao/go-bootcamp/internal/job"
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
	InitSaramaClient,
	InitSyncProducer,
)

var interactiveSvcSet = wire.NewSet(
	intrDao.NewGORMInteractiveDAO,
	intrRediscache.NewInteractiveRedisCache,
	ioc.InitTopArticlesCache,
	intrRepository.NewCachedInteractiveRepository,
	intrService.NewInteractiveService,
)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptJobRepository,
	dao.NewGORMJobDAO,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// third-party dependencies
		thirdPartySet,

		interactiveSvcSet,
		ioc.InitIntrClient,

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

		article.NewSaramaSyncProducer,
		// article.NewInteractiveReadEventConsumer,

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

func InitArticleHandler(articleDAO dao.ArticleDAO) *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		rediscache.NewArticleRedisCache,
		interactiveSvcSet,
		ioc.InitIntrClient,

		rediscache.NewUserRedisCache,
		dao.NewUserDAO,
		repository.NewUserRepository,

		article.NewSaramaSyncProducer,

		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}

func InitJobScheduler() *job.Scheduler {
	wire.Build(jobProviderSet, thirdPartySet, job.NewScheduler)
	return &job.Scheduler{}
}
