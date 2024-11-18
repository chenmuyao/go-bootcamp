package ratelimit

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	RedisClient *redis.Client
	Interval    time.Duration
	Limit       int
}

func NewBuilder(redisClient *redis.Client, interval time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		RedisClient: redisClient,
		Interval:    interval,
		Limit:       limit,
	}
}

func (rl *RateLimiter) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
