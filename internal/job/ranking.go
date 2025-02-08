package job

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/service"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

// Name implements Job.
func (r *RankingJob) Name() string {
	return "ranking"
}

// Run implements Job.
func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) Job {
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
	}
}
