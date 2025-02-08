package job

import (
	"strconv"
	"time"

	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	l      logger.Logger
	vector *prometheus.SummaryVec
}

func NewCronJobBuilder(l logger.Logger, opts prometheus.SummaryOpts) *CronJobBuilder {
	vector := promauto.NewSummaryVec(opts, []string{"name", "success"})
	return &CronJobBuilder{
		l:      l,
		vector: vector,
	}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobAdapterFunc(func() {
		start := time.Now()
		c.l.Debug("start job", logger.String("name", name))
		err := job.Run()
		if err != nil {
			c.l.Error("Failed to execute job", logger.String("name", name), logger.Error(err))
		}
		c.l.Debug("stop job", logger.String("name", name))
		duration := time.Since(start).Milliseconds()
		c.vector.WithLabelValues(name, strconv.FormatBool(err == nil)).Observe(float64(duration))
	})
}

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}
