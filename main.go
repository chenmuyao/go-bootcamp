package main

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := gin.Default()

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

	db, err := gorm.Open(
		mysql.Open(
			"root:root@tcp(127.0.0.1:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local",
		),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}

	dao.InitTable(db)

	userDAO := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDAO)
	userService := service.NewUserService(userRepo)
	userHandlers := web.NewUserHandler(userService)
	userHandlers.RegisterRoutes(server)

	server.Run(":7779")
}
