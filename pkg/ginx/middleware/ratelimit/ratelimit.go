package ratelimit

import (
	"net/http"

	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	FixedWindowLimiter   limiter.FixedWindowLimiter
	SlidingWindowLimiter limiter.SlidingWindowLimiter
	limiterType          int
}

type FixedWindowOptions limiter.FixedWindowOptions

func NewFixedWindowLimiterBuilder(options *FixedWindowOptions) *RateLimiter {
	return &RateLimiter{
		FixedWindowLimiter: *limiter.NewFixedWindowLimiter(&limiter.FixedWindowOptions{
			Interval: options.Interval,
			Limit:    options.Limit,
			Prefix:   options.Prefix,
		}),
		limiterType: limiter.FixedWindow,
	}
}

type SlidingWindowOptions limiter.SlidingWindowOptions

func NewSlidingWindowLimiterBuilder(options *SlidingWindowOptions) *RateLimiter {
	return &RateLimiter{
		SlidingWindowLimiter: *limiter.NewSlidingWindowLimiter(&limiter.SlidingWindowOptions{
			Interval:      options.Interval,
			WindowsAmount: options.WindowsAmount,
			Limit:         options.Limit,
			Prefix:        options.Prefix,
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
