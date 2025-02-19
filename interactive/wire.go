//go:build wireinject
// +build wireinject

package main

import (
	"github.com/chenmuyao/go-bootcamp/interactive/events"
	"github.com/chenmuyao/go-bootcamp/interactive/grpc"
	"github.com/chenmuyao/go-bootcamp/interactive/ioc"
	intrRepository "github.com/chenmuyao/go-bootcamp/interactive/repository"
	intrRediscache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache/rediscache"
	intrDao "github.com/chenmuyao/go-bootcamp/interactive/repository/dao"
	intrService "github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitRedis,
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
)

var interactiveSvcSet = wire.NewSet(
	intrDao.NewGORMInteractiveDAO,
	intrRediscache.NewInteractiveRedisCache,
	ioc.InitTopArticlesCache,
	intrRepository.NewCachedInteractiveRepository,
	intrService.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}

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
