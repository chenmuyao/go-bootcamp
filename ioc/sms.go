package ioc

import (
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/localsms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(redisClient redis.Cmdable) sms.Service {
	return ratelimit.NewRateLimitSMSService(
		localsms.NewService(),
		limiter.NewLimiter(&limiter.RedisTokenBucketOptions{
			RedisClient:   redisClient,
			Prefix:        "",
			Capacity:      100,
			ReleaseAmount: 10,
			Interval:      1 * time.Second,
		}),
	)
	// return &localsms.Service{}
}
