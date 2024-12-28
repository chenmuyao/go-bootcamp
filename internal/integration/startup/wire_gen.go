// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache/rediscache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	cmdable := InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	logger := InitLogger()
	v := ioc.InitGinMiddlewares(cmdable, handler, logger)
	db := InitDB()
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
	oAuth2GiteaHandler := NewDummyGiteaHandler(userService, handler, logger)
	articleDAO := dao.NewArticleDAO(db)
	articleCache := rediscache.NewArticleRedisCache(cmdable)
	articleRepository := repository.NewArticleRepository(logger, articleDAO, articleCache, userRepository)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(logger, articleService)
	engine := ioc.InitWebServer(v, userHandler, oAuth2GiteaHandler, articleHandler)
	return engine
}

func InitArticleHandler(articleDAO dao.ArticleDAO) *web.ArticleHandler {
	logger := InitLogger()
	cmdable := InitRedis()
	articleCache := rediscache.NewArticleRedisCache(cmdable)
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userCache := rediscache.NewUserRedisCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	articleRepository := repository.NewArticleRepository(logger, articleDAO, articleCache, userRepository)
	articleService := service.NewArticleService(articleRepository)
	articleHandler := web.NewArticleHandler(logger, articleService)
	return articleHandler
}

// wire.go:

var thirdPartySet = wire.NewSet(
	InitRedis,
	InitDB,
	InitLogger,
)
