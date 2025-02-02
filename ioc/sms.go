package ioc

import (
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/localsms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/opentelemetry"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/prometheus"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
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

	promSvc := prometheus.NewPrometheusSMS(rateLimitSMSSvc, prom.SummaryOpts{
		Namespace: "my_company",
		Subsystem: "wetravel",
		Name:      "sms_svc",
		Help:      "SMS service metrics",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
		ConstLabels: prom.Labels{
			"instance_id": "instance",
		},
	})

	tracer := otel.Tracer("github.com/blabla/opentelemetry")
	otelSvc := opentelemetry.NewOTELSMSSvc(promSvc, tracer)
	return otelSvc
}
