package failover

import (
	"context"
	"sync/atomic"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

type TimeoutFailOverSMSService struct {
	svcs []sms.Service

	// Current index
	idx int32

	// TO counter
	cnt int32

	// read-only
	threshold int32
}

func NewTimeoutFailOverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

// }}}
// {{{ Other structs

// }}}
// {{{ Struct Methods

func (t *TimeoutFailOverSMSService) Send(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		// if index is already updated by another thread
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// reset cnt
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = newIdx
	}

	// NOTE: Not strict N TO, could have concurrency issue, but it is not
	// fatal for this feature.

	svc := t.svcs[idx]
	err := svc.Send(ctx, toNb, body, args...)
	switch err {
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
	default:
		// If strict TO switch, keep cnt unchanged
		// If N errors, cnt ++
		// If EOF --> change svc directly.
	}
	return err
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
