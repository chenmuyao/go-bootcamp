package ioc

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/web"
	"github.com/chenmuyao/go-bootcamp/internal/web/middleware"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx/middleware/ratelimit"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func InitWebServer(middlewares []gin.HandlerFunc, userHandlers *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandlers.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
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
		}),

		useJWT(),
		// useSession(),
		// sessionCheckLogin(),

		// ratelimit.NewFixedWindowLimiterBuilder(&ratelimit.FixedWindowOptions{
		// 	Interval: time.Second,
		// 	Limit:    100,
		// }).Build(),
		ratelimit.NewSlidingWindowLimiterBuilder(&ratelimit.SlidingWindowOptions{
			WindowSize: time.Second,
			Limit:      100,
		}).Build(),
	}
}

func useJWT() gin.HandlerFunc {
	loginJWT := middleware.NewLoginJWT([]string{
		"/user/signup",
		"/user/login",
		"/user/login_sms/code/send",
		"/user/login_sms",
	})
	return loginJWT.CheckLogin()
}

func sessionCheckLogin() gin.HandlerFunc {
	login := middleware.LoginMiddleware([]string{
		"/user/signup",
		"/user/login",
		"/user/login_sms/code/send",
		"/user/login_sms",
	})
	return login.CheckLogin()
}

func useSession() gin.HandlerFunc {
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
	return sessions.Sessions("ssid", store)
}
