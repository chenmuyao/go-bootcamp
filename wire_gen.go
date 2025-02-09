// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/chenmuyao/go-bootcamp/internal/events/article"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
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
	articleDAO := dao.NewArticleDAO(db)
	articleCache := rediscache.NewArticleRedisCache(cmdable)
	userDAO := dao.NewUserDAO(db)
	userCache := rediscache.NewUserRedisCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	articleRepository := repository.NewArticleRepository(logger, articleDAO, articleCache, userRepository)
	interactiveRepository := repository.NewCachedInteractiveRepository(logger, interactiveDAO, interactiveCache, topArticlesCache, articleRepository)
	return interactiveRepository
}

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := ioc.InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := ioc.InitDB(logger)
	userDAO := dao.NewUserDAO(db)
	userCache := rediscache.NewUserRedisCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := rediscache.NewCodeRedisCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	asyncSMSDAO := dao.NewAsyncSMSDAO(db)
	asyncSMSRepository := repository.NewAsyncSMSRepository(asyncSMSDAO, db)
	smsService := ioc.InitSMSService(cmdable, asyncSMSRepository)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(logger, userService, codeService, handler)
	giteaService := ioc.InitGiteaService(logger)
	oAuth2GiteaHandler := web.NewOAuth2GiteaHandler(logger, giteaService, userService, handler)
	articleDAO := dao.NewArticleDAO(db)
	articleCache := rediscache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewArticleRepository(logger, articleDAO, articleCache, userRepository)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitSyncProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, producer)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	interactiveCache := rediscache.NewInteractiveRedisCache(cmdable)
	topArticlesCache := ioc.InitTopArticlesCache()
	interactiveRepository := repository.NewCachedInteractiveRepository(logger, interactiveDAO, interactiveCache, topArticlesCache, articleRepository)
	interactiveService := service.NewInteractiveService(interactiveRepository)
	articleHandler := web.NewArticleHandler(logger, articleService, interactiveService)
	engine := ioc.InitWebServer(v, userHandler, oAuth2GiteaHandler, articleHandler)
	interactiveReadEventConsumer := article.NewInteractiveReadEventConsumer(logger, interactiveRepository, client)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer)
	rankingCache := rediscache.NewRankingRedisCache(cmdable)
	rankingLocalCache := ioc.InitRankingLocalCache()
	rankingRepository := repository.NewCachedRankingRepository(rankingCache, rankingLocalCache)
	rankingService := service.NewBatchRankingService(interactiveService, articleService, rankingRepository)
	job := ioc.InitRankingJob(rankingService, logger, cmdable)
	cron := ioc.InitJobs(logger, job)
	app := &App{
		server:    engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

// wire.go:

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, rediscache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)

var rankingSvcSet = wire.NewSet(ioc.InitRankingLocalCache, rediscache.NewRankingRedisCache, repository.NewCachedRankingRepository, service.NewBatchRankingService)
