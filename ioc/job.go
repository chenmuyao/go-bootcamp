package ioc

import (
	"time"

	"github.com/bsm/redislock"
	"github.com/chenmuyao/go-bootcamp/internal/job"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService, l logger.Logger, redis redis.Cmdable) job.Job {
	lock := redislock.New(redis)
	return job.NewRankingJob(svc, lock, time.Second*30, l)
}

func InitJobs(l logger.Logger, j job.Job) *cron.Cron {
	builder := job.NewCronJobBuilder(l, prometheus.SummaryOpts{
		Namespace: "my_company",
		Subsystem: "wetravel",
		Name:      "cron_job",
		Help:      "cron job",
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("@every 1m", builder.Build(j))
	if err != nil {
		panic(err)
	}
	return expr
}
