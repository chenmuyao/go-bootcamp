package main

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	userHandlers := web.NewUserHandler()
	userHandlers.RegisterRoutes(server)

	server.Run(":7779")
}
