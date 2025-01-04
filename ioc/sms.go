package ioc

import (
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/localsms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	"github.com/redis/go-redis/v9"
)

func InitSMSService(
	redisClient redis.Cmdable,
	asyncRepo repository.AsyncSMSRepository,
) sms.Service {
	rateLimitSMSSvc := ratelimit.NewRateLimitSMSService(
		localsms.NewService(),
		limiter.NewLimiter(&limiter.RedisTokenBucketOptions{
			RedisClient:   redisClient,
			Prefix:        "sms-svc",
			Capacity:      100,
			ReleaseAmount: 10,
			Interval:      1 * time.Second,
		}),
	)
	// TODO: replace the context by a global shutdown context
	//    asyncSvc := async.NewAsyncSMSService(
	// 	context.Background(),
	// 	rateLimitSMSSvc,
	// 	asyncRepo,
	// 	&async.AsyncSMSServiceOptions{
	// 		PollInterval:    10 * time.Second,
	// 		RetryTimes:      3,
	// 		RetryErrorCodes: []int{20504},
	// 	},
	// )
	return rateLimitSMSSvc
}
