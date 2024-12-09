package ratelimit

import (
	"net/http"

	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	limiter limiter.Limiter
}

func NewRateLimiterBuilder(options any) *RateLimiter {
	l := limiter.NewLimiter(options)
	return &RateLimiter{
		limiter: l,
	}
}

func (rl *RateLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if !rl.limiter.AcceptConnection(ctx, ip) {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
