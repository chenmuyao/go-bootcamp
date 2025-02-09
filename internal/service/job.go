package service

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
)

type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
	// NOTE: method to release the lock, but we prefer using a cancelfunc embedded in job
	// Release(ctx context.Context, job domain.Job) error
}

type cronJobService struct {
	l               logger.Logger
	repo            repository.JobRepository
	refreshInterval time.Duration
}

// ResetNextTime implements JobService.
func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	nextTime := job.NextTime()
	return c.repo.UpdateNextTime(ctx, job.ID, nextTime)
}

// Preempt implements JobService.
func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			c.refresh(job.ID)
		}
	}()
	job.CancelFunc = func() {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.repo.Release(ctx, job.ID)
		if err != nil {
			c.l.Error("failed to release the job", logger.Int64("jid", job.ID), logger.Error(err))
		}
	}
	return job, nil
}

func (c *cronJobService) refresh(jid int64) {
	// update utime
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, jid)
	if err != nil {
		c.l.Error("failed to refresh", logger.Int64("jid", jid), logger.Error(err))
	}
}

func NewCronJobService(
	l logger.Logger,
	repo repository.JobRepository,
) JobService {
	return &cronJobService{
		l:               nil,
		repo:            repo,
		refreshInterval: time.Minute,
	}
}
