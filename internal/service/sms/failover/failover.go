package failover

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
)

var errFailOverAll = errors.New("failed after polling all service providers")

type FailOverSMSService struct {
	svcs []sms.Service

	// V1: dynamically compute the polling start point
	idx uint64
}

func NewFailOverSMSService(svcs []sms.Service) *FailOverSMSService {
	return &FailOverSMSService{
		svcs: svcs,
	}
}

func (f *FailOverSMSService) Send(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, toNb, body, args...)
		if err == nil {
			return nil
		}
		slog.Error("failover send error", "err", err)
	}
	return errFailOverAll
}

// Polling from the last position. Polling all services equally.
// Exclude user canceling and time out.
func (f *FailOverSMSService) SendV1(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	// NOTE: not atomic, different CPUs can read in different idx from their
	// CPU cache
	// idx := f.idx + 1
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[i&length]
		err := svc.Send(ctx, toNb, body, args...)
		switch err {
		case nil:
			return nil
		case context.Canceled, context.DeadlineExceeded:
			// canceled by user or timed out
			return err
		}
		slog.Error("failover send error", "err", err)
	}
	return errFailOverAll
}
