package main

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := initWebServer()

	db := initDB()

	initUserHandlers(db, server)

	server.Run(":7779")
}

func initDB() *gorm.DB {
	db, err := gorm.Open(
		mysql.Open(
			"root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
		),
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

	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: "localhost:6379",
	// })

	// server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	useJWT(server)

	return server
}

func initUserHandlers(db *gorm.DB, server *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	userHandlers := web.NewUserHandler(userService)
	userHandlers.RegisterRoutes(server)
}

func useJWT(server *gin.Engine) {
	loginJWT := middleware.LoginJWTMiddleware([]string{
		"/user/signup",
		"/user/login",
	})
	server.Use(loginJWT.CheckLogin())
}

func useSession(server *gin.Engine) {
	login := middleware.LoginMiddleware([]string{
		"/user/signup",
		"/user/login",
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
