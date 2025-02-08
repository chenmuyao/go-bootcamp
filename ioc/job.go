package ioc

import (
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/job"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService) job.Job {
	return job.NewRankingJob(svc, time.Second*30)
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
