package main

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/middleware"
	"github.com/chenmuyao/go-bootcamp/internal/web/validate"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	validate.UseValidators("date")

	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "my_company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	login := middleware.LoginMiddleware([]string{
		"/user/signup",
		"/user/login",
	})

	// create store to hold sessions in Cookies
	store := cookie.NewStore([]byte("secret"))

	// Use the store to hold session ssid
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())

	return server
}

func initUserHandlers(db *gorm.DB, server *gin.Engine) {
	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	userHandlers := web.NewUserHandler(userService)
	userHandlers.RegisterRoutes(server)
}
