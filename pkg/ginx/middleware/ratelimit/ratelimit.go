package ratelimit

import (
	"net/http"
	"time"

	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	FixedWindowLimiter   limiter.FixedWindowLimiter
	SlidingWindowLimiter limiter.SlidingWindowLimiter
	limiterType          int
}

type FixedWindowOptions struct {
	Limit    int
	Interval time.Duration
}

func NewFixedWindowLimiterBuilder(options *FixedWindowOptions) *RateLimiter {
	return &RateLimiter{
		FixedWindowLimiter: *limiter.NewFixedWindowLimiter(&limiter.FixedWindowOptions{
			Interval: options.Interval,
			Limit:    options.Limit,
		}),
		limiterType: limiter.FixedWindow,
	}
}

type SlidingWindowOptions struct {
	Limit      int
	WindowSize time.Duration
}

func NewSlidingWindowLimiterBuilder(options *SlidingWindowOptions) *RateLimiter {
	return &RateLimiter{
		SlidingWindowLimiter: *limiter.NewSlidingWindowLimiter(&limiter.SlidingWindowOptions{
			WindowSize: options.WindowSize,
			Limit:      options.Limit,
		}),
		limiterType: limiter.SlidingWindow,
	}
}

func (rl *RateLimiter) Build() gin.HandlerFunc {
	var acceptConnection func(string) bool
	switch rl.limiterType {
	case limiter.FixedWindow:
		acceptConnection = rl.FixedWindowLimiter.AcceptConnection
	case limiter.SlidingWindow:
		acceptConnection = rl.SlidingWindowLimiter.AcceptConnection
	}

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if !acceptConnection(ip) {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}
