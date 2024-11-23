package main

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/config"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/cache"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/localsms"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/middleware"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := initWebServer()

	db := initDB()
	cache := initCache()

	codeService := initCodeSvc(cache)
	initUserHandlers(db, cache, codeService, server)

	server.Run(":8081")
}

func initCache() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
}

func initDB() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(config.Config.DB.DSN),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}

	err = dao.InitTable(db)
	if err != nil {
		panic("failed to init tables")
	}
	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		// Allow frontend to access headers sent back from the backend
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "my_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	// server.Use(ratelimit.NewFixedWindowLimiterBuilder(&ratelimit.FixedWindowOptions{
	// 	Interval: time.Second,
	// 	Limit:    100,
	// }).Build())
	server.Use(ratelimit.NewSlidingWindowLimiterBuilder(&ratelimit.SlidingWindowOptions{
		WindowSize: time.Second,
		Limit:      100,
	}).Build())

	useJWT(server)

	return server
}

func initUserHandlers(
	db *gorm.DB,
	redisClient redis.Cmdable,
	codeService *service.CodeService,
	server *gin.Engine,
) {
	userDAO := dao.NewUserDAO(db)
	userCache := cache.NewUserCache(redisClient)
	userRepo := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepo)
	userHandlers := web.NewUserHandler(userService, codeService)
	userHandlers.RegisterRoutes(server)
}

func initCodeSvc(redisClient redis.Cmdable) *service.CodeService {
	sms := localsms.NewService()
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	return service.NewCodeService(crepo, sms)
}

func useJWT(server *gin.Engine) {
	loginJWT := middleware.NewLoginJWT([]string{
		"/user/signup",
		"/user/login",
		"/user/login_sms/code/send",
		"/user/login_sms",
	})
	server.Use(loginJWT.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := middleware.LoginMiddleware([]string{
		"/user/signup",
		"/user/login",
		"/user/login_sms/code/send",
		"/user/login_sms",
	})

	// create store to hold session data in Cookies
	// store := cookie.NewStore([]byte("secret"))

	// create store to hold session data in memstore
	// store := memstore.NewStore(
	// 	[]byte("QbYQn3ZyECBq3fQwWFj84ccoqipj70oJ"),
	// 	[]byte("kpqqi5guoJGKCmsgN7a5jwgd2nvpC2P3"),
	// )

	// NOTE: Use redis for distributed storage of session info
	store, err := redisStore.NewStore(
		16,
		"tcp",
		"localhost:6379",
		"",
		[]byte("QbYQn3ZyECBq3fQwWFj84ccoqipj70oJ"), // authentication
		[]byte("kpqqi5guoJGKCmsgN7a5jwgd2nvpC2P3"), // encryption
	)
	if err != nil {
		panic(err)
	}

	// Use the store to hold session ssid
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
