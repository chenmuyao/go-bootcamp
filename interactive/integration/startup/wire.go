//go:build wireinject
// +build wireinject

package startup

import (
	"github.com/chenmuyao/go-bootcamp/interactive/grpc"
	intrRepository "github.com/chenmuyao/go-bootcamp/interactive/repository"
	intrRediscache "github.com/chenmuyao/go-bootcamp/interactive/repository/cache/rediscache"
	intrDao "github.com/chenmuyao/go-bootcamp/interactive/repository/dao"
	intrService "github.com/chenmuyao/go-bootcamp/interactive/service"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	InitRedis,
	InitDB,
	InitLogger,
	InitSaramaClient,
)

var interactiveSvcSet = wire.NewSet(
	intrDao.NewGORMInteractiveDAO,
	intrRediscache.NewInteractiveRedisCache,
	ioc.InitTopArticlesCache,
	intrRepository.NewCachedInteractiveRepository,
	intrService.NewInteractiveService,
)

func InitInteractiveService() *grpc.InteractiveServiceServer {
	wire.Build(thirdPartySet, interactiveSvcSet, grpc.NewInteractiveServiceServer)
	return &grpc.InteractiveServiceServer{}
}
