package repository

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/repository/dao"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jid int64) error
	UpdateUtime(ctx context.Context, jid int64) error
	UpdateNextTime(ctx context.Context, jid int64, nextTime time.Time) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

// UpdateNextTime implements JobRepository.
func (p *PreemptJobRepository) UpdateNextTime(
	ctx context.Context,
	jid int64,
	nextTime time.Time,
) error {
	return p.dao.UpdateNextTime(ctx, jid, nextTime)
}

// UpdateUtime implements JobRepository.
func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, jid int64) error {
	return p.dao.UpdateUtime(ctx, jid)
}

// Preempt implements CronJobRepository.
func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		ID:       j.ID,
		CronExpr: j.CronExpr,
		Executor: j.Executor,
		Name:     j.Name,
	}, nil
}

// Release implements CronJobRepository.
func (p *PreemptJobRepository) Release(ctx context.Context, jid int64) error {
	return p.dao.Release(ctx, jid)
}

func NewPreemptJobRepository(dao dao.JobDAO) JobRepository {
	return &PreemptJobRepository{
		dao: dao,
	}
}
