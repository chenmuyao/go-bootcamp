// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/chenmuyao/go-bootcamp/interactive/events"
	"github.com/chenmuyao/go-bootcamp/interactive/repository"
	"github.com/chenmuyao/go-bootcamp/interactive/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/interactive/repository/dao"
	service2 "github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/internal/events/article"
	"github.com/chenmuyao/go-bootcamp/internal/job"
	repository2 "github.com/chenmuyao/go-bootcamp/internal/repository"
	rediscache2 "github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	dao2 "github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitInteractiveRepo() repository.InteractiveRepository {
	logger := ioc.InitLogger()
	db := ioc.InitDB(logger)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := rediscache.NewInteractiveRedisCache(cmdable)
	topArticlesCache := ioc.InitTopArticlesCache()
	interactiveRepository := repository.NewCachedInteractiveRepository(logger, interactiveDAO, interactiveCache, topArticlesCache)
	return interactiveRepository
}

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := ioc.InitDB(logger)
	userDAO := dao2.NewUserDAO(db)
	userCache := rediscache2.NewUserRedisCache(cmdable)
	userRepository := repository2.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := rediscache2.NewCodeRedisCache(cmdable)
	codeRepository := repository2.NewCodeRepository(codeCache)
	asyncSMSDAO := dao2.NewAsyncSMSDAO(db)
	asyncSMSRepository := repository2.NewAsyncSMSRepository(asyncSMSDAO, db)
	smsService := ioc.InitSMSService(cmdable, asyncSMSRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(logger, userService, codeService, handler)
	giteaService := ioc.InitGiteaService(logger)
	oAuth2GiteaHandler := web.NewOAuth2GiteaHandler(logger, giteaService, userService, handler)
	articleDAO := dao2.NewArticleDAO(db)
	articleCache := rediscache2.NewArticleRedisCache(cmdable)
	articleRepository := repository2.NewArticleRepository(logger, articleDAO, articleCache, userRepository)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := rediscache.NewInteractiveRedisCache(cmdable)
	topArticlesCache := ioc.InitTopArticlesCache()
	interactiveRepository := repository.NewCachedInteractiveRepository(logger, interactiveDAO, interactiveCache, topArticlesCache)
	interactiveService := service2.NewInteractiveService(interactiveRepository)
	interactiveServiceClient := ioc.InitIntrClient(interactiveService)
	articleHandler := web.NewArticleHandler(logger, articleService, interactiveServiceClient)
	engine := ioc.InitWebServer(v, userHandler, oAuth2GiteaHandler, articleHandler)
	interactiveReadEventConsumer := events.NewInteractiveReadEventConsumer(logger, interactiveRepository, client)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	rankingCache := rediscache2.NewRankingRedisCache(cmdable)
	rankingLocalCache := ioc.InitRankingLocalCache()
	rankingRepository := repository2.NewCachedRankingRepository(rankingCache, rankingLocalCache)
	rankingService := service.NewBatchRankingService(interactiveServiceClient, articleService, rankingRepository)
	job := ioc.InitRankingJob(rankingService, logger, cmdable)
	cron := ioc.InitJobs(logger, job)
	app := &App{
		server:    engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

func InitJobScheduler() *job.Scheduler {
	logger := ioc.InitLogger()
	db := ioc.InitDB(logger)
	jobDAO := dao2.NewGORMJobDAO(db, logger)
	jobRepository := repository2.NewPreemptJobRepository(jobDAO)
	jobService := service.NewCronJobService(logger, jobRepository)
	scheduler := job.NewScheduler(logger, jobService)
	return scheduler
}

// wire.go:

var thirdPartySet = wire.NewSet(ioc.InitRedis, ioc.InitDB, ioc.InitLogger, ioc.InitSaramaClient, ioc.InitSyncProducer)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, rediscache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service2.NewInteractiveService)

var rankingSvcSet = wire.NewSet(ioc.InitRankingLocalCache, rediscache2.NewRankingRedisCache, repository2.NewCachedRankingRepository, service.NewBatchRankingService)

var jobProviderSet = wire.NewSet(service.NewCronJobService, repository2.NewPreemptJobRepository, dao2.NewGORMJobDAO)
