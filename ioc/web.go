package ioc

import (
	"strings"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/web"
	ijwt "github.com/chenmuyao/go-bootcamp/internal/web/jwt"
	"github.com/chenmuyao/go-bootcamp/internal/web/middleware"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx"
	"github.com/chenmuyao/go-bootcamp/pkg/ginx/middleware/prometheus"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func InitWebServer(middlewares []gin.HandlerFunc,
	userHandlers *web.UserHandler,
	giteaHandlers *web.OAuth2GiteaHandler,
	articleHandlers *web.ArticleHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	userHandlers.RegisterRoutes(server)
	giteaHandlers.RegisterRoutes(server)
	articleHandlers.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(
	redisClient redis.Cmdable,
	jwtHdl ijwt.Handler,
	l logger.Logger,
) []gin.HandlerFunc {
	pb := prometheus.NewPrometheusBuilder(
		"my_company",
		"wetravel",
		"gin_http",
		"instance",
		"Gin requests data",
	)
	ginx.InitCounter(prom.CounterOpts{
		Namespace: "my_company",
		Subsystem: "wetravel",
		Name:      "errcode",
		Help:      "Error code data",
		ConstLabels: prom.Labels{
			"instance_id": "instance",
		},
	})
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			// Allow frontend to access headers sent back from the backend
			ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					return true
				}
				return strings.Contains(origin, "my_company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		pb.BuildResponseTime(),
		pb.BuildActiveRequest(),
		otelgin.Middleware("wetravel"),

		// ratelimit.NewRateLimiterBuilder(&limiter.RedisSlidingWindowOptions{
		// 	Prefix:        "web-interface",
		// 	RedisClient:   redisClient,
		// 	Interval:      time.Second,
		// 	Limit:         100,
		// 	WindowsAmount: 10,
		// }).Build(),
		// middleware.NewLogMiddlewareBuilder(func(ctx context.Context, al middleware.AccessLog) {
		// 	l.Debug("", logger.Field{Key: "req", Value: al})
		// }).AllowReqBody().AllowRespBody().Build(),
		useJWT(jwtHdl),
		// useSession(),
		// sessionCheckLogin(),

	}
}

func useJWT(jwtHdl ijwt.Handler) gin.HandlerFunc {
	loginJWT := middleware.NewLoginJWTBuilder(jwtHdl, []string{
		"/user/signup",
		"/user/login",
		"/user/login_sms/code/send",
		"/user/login_sms",
		"/user/refresh_token",
		"/oauth2/gitea/authurl",
		"/oauth2/gitea/callback",
	})
	return loginJWT.Build()
}

// func sessionCheckLogin() gin.HandlerFunc {
// 	login := middleware.LoginMiddleware([]string{
// 		"/user/signup",
// 		"/user/login",
// 		"/user/login_sms/code/send",
// 		"/user/login_sms",
// 	})
// 	return login.CheckLogin()
// }

// func useSession() gin.HandlerFunc {
// 	// create store to hold session data in Cookies
// 	// store := cookie.NewStore([]byte("secret"))
//
// 	// create store to hold session data in memstore
// 	// store := memstore.NewStore(
// 	// 	[]byte("QbYQn3ZyECBq3fQwWFj84ccoqipj70oJ"),
// 	// 	[]byte("kpqqi5guoJGKCmsgN7a5jwgd2nvpC2P3"),
// 	// )
//
// 	// NOTE: Use redis for distributed storage of session info
// 	store, err := redisStore.NewStore(
// 		16,
// 		"tcp",
// 		"localhost:6379",
// 		"",
// 		[]byte("QbYQn3ZyECBq3fQwWFj84ccoqipj70oJ"), // authentication
// 		[]byte("kpqqi5guoJGKCmsgN7a5jwgd2nvpC2P3"), // encryption
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// Use the store to hold session ssid
// 	return sessions.Sessions("ssid", store)
// }
