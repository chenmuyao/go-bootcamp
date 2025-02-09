package job

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

type RankingJob struct {
	l          logger.Logger
	svc        service.RankingService
	timeout    time.Duration
	lockClient *redislock.Client
}

// Name implements Job.
func (r *RankingJob) Name() string {
	return "ranking"
}

// Run implements Job.
func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()
	lock, err := r.lockClient.Obtain(ctx, "job:ranking", r.timeout, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3),
	})
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		er := lock.Release(ctx)
		if er != nil {
			r.l.Error("ranking job failed to release distributed loc", logger.Error(er))
		}
	}()

	bizCtx, bizCancel := context.WithTimeout(context.Background(), r.timeout)
	defer bizCancel()
	return r.svc.TopN(bizCtx)
}

func NewRankingJob(
	svc service.RankingService,
	lock *redislock.Client,
	timeout time.Duration,
	l logger.Logger,
) Job {
	return &RankingJob{
		l:          l,
		svc:        svc,
		timeout:    timeout,
		lockClient: lock,
	}
}
