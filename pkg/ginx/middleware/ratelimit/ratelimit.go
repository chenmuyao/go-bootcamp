package ratelimit

import (
	"net/http"
	"time"

	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Options struct {
	RedisClient *redis.Client
	Interval    time.Duration
	Limit       int
}

type RateLimiter struct {
	FixedWindowLimiter limiter.FixedWindowLimiter
}

func NewBuilder(options *Options) *RateLimiter {
	return &RateLimiter{
		FixedWindowLimiter: *limiter.NewFixedWindowLimiter(&limiter.Options{
			Interval: options.Interval,
			Limit:    options.Limit,
		}),
	}
}

func (rl *RateLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if !rl.FixedWindowLimiter.AcceptConnection(ip) {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
