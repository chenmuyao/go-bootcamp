package job

import (
	"context"
	"fmt"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service"
	"github.com/chenmuyao/go-bootcamp/pkg/logger"
	"golang.org/x/sync/semaphore"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
	RegisterExecutor(name string, fn func(ctx context.Context, j domain.Job) error)
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

// Exec implements Executor.
func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("func %s unregistered", j.Name)
	}
	return fn(ctx, j)
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func (l *LocalFuncExecutor) RegisterExecutor(
	name string,
	fn func(ctx context.Context, j domain.Job) error,
) {
	l.funcs[name] = fn
}

func NewLocalFuncExecutor() Executor {
	return &LocalFuncExecutor{funcs: map[string]func(ctx context.Context, j domain.Job) error{}}
}

type Scheduler struct {
	l         logger.Logger
	dbTimeout time.Duration
	svc       service.JobService

	executors map[string]Executor
	limiter   semaphore.Weighted
}

func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.executors[exec.Name()] = exec
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		// Check if we should run the scheduler at all
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		dbCtx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		exec, ok := s.executors[j.Executor]
		if !ok {
			s.l.Error(
				"failed to find executor",
				logger.Int64("jid", j.ID),
				logger.String("executor", j.Executor),
			)
			continue
		}

		go func() {
			defer func() {
				s.limiter.Release(1)
				j.CancelFunc()
			}()

			er := exec.Exec(ctx, j)
			if er != nil {
				s.l.Error("failed to run job", logger.Int64("jid", j.ID), logger.Error(err))
				return
			}

			er = s.svc.ResetNextTime(ctx, j)
			if er != nil {
				s.l.Error("failed to reset next time", logger.Int64("jid", j.ID), logger.Error(err))
			}
		}()
	}
}

func NewScheduler(
	l logger.Logger,
	svc service.JobService,
) *Scheduler {
	return &Scheduler{
		l:         l,
		dbTimeout: time.Second,
		svc:       svc,
		executors: map[string]Executor{},
		limiter:   *semaphore.NewWeighted(100),
	}
}
