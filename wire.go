//go:build wireinject
// +build wireinject

package main

import (
	intrEvents "github.com/chenmuyao/go-bootcamp/interactive/events"
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
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitSyncProducer,
)

var interactiveSvcSet = wire.NewSet(
	intrDao.NewGORMInteractiveDAO,
	intrRediscache.NewInteractiveRedisCache,
	intrRepository.NewCachedInteractiveRepository,
	intrService.NewInteractiveService,
)

var rankingSvcSet = wire.NewSet(
	ioc.InitRankingLocalCache,
	rediscache.NewRankingRedisCache,
	repository.NewCachedRankingRepository,
	service.NewBatchRankingService,
)

var jobProviderSet = wire.NewSet(
	service.NewCronJobService,
	repository.NewPreemptJobRepository,
	dao.NewGORMJobDAO,
)

func InitInteractiveRepo() intrRepository.InteractiveRepository {
	wire.Build(
		thirdPartySet,

		ioc.InitTopArticlesCache,

		intrDao.NewGORMInteractiveDAO,

		intrRediscache.NewInteractiveRedisCache,

		intrRepository.NewCachedInteractiveRepository,
	)
	return &intrRepository.CachedInteractiveRepository{}
}

func InitWebServer() *App {
	wire.Build(
		// third-party dependencies
		thirdPartySet,

		interactiveSvcSet,
		rankingSvcSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		article.NewSaramaSyncProducer,
		intrEvents.NewInteractiveReadEventConsumer,
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

func InitJobScheduler() *job.Scheduler {
	wire.Build(jobProviderSet, thirdPartySet, job.NewScheduler)
	return &job.Scheduler{}
}
