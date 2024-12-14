package async

import (
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/repository"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/internal/service/sms/ratelimit"
)

// {{{ Consts

// }}}
// {{{ Global Varirables

// }}}
// {{{ Interface

// }}}
// {{{ Struct

// NOTE: Use Singleton to fire only one go routine that reprogram async send.
type AsyncSMSService struct {
	repo         repository.AsyncSMSRepository
	svc          sms.Service
	pollInterval time.Duration
	retryTimes   int
	// For twilio, it could be 20504 Twilio Internal Error
	retryErrorCodes []int
}

func NewAsyncSMSService(
	ctx context.Context,
	smsSvc sms.Service,
	asyncRepo repository.AsyncSMSRepository,
	opts *AsyncSMSServiceOptions,
) *AsyncSMSService {
	a := &AsyncSMSService{
		repo:            asyncRepo,
		svc:             smsSvc,
		retryTimes:      opts.RetryTimes,
		retryErrorCodes: opts.RetryErrorCodes,
		pollInterval:    opts.PollInterval,
	}
	go a.asyncSend(ctx)
	return a
}

// }}}
// {{{ Other structs

type AsyncSMSServiceOptions struct {
	PollInterval time.Duration
	RetryTimes   int
	// For twilio, it could be 20504 Twilio Internal Error
	RetryErrorCodes []int
}

// }}}
// {{{ Struct Methods

func (a *AsyncSMSService) Send(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	err := a.svc.Send(ctx, toNb, body, args...)

	// Retry for certain error codes
	if err, ok := err.(*sms.SMSThirdPartyError); ok {
		if slices.Contains(a.retryErrorCodes, err.Code) {
			slog.Error("async send sms because of third party errcode", "err", err)
			return a.store(ctx, toNb, body, args...)
		}
		return err
	}

	switch err {
	case nil:
		return nil
	case ratelimit.ErrLimited:
		slog.Error("async send sms because of rate limit", "err", err)
		return a.store(ctx, toNb, body, args...)
	default:
		return err
	}
}

func (a *AsyncSMSService) asyncSend(ctx context.Context) {
	// TODO: graceful shutdown to be implemented globally and pass a shutdown context here.
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a.repo.TrySend(ctx, a.Send, a.retryTimes)
			time.Sleep(a.pollInterval)
		}
	}
}

func (a *AsyncSMSService) store(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return err
	}

	return a.repo.AddSMS(ctx, toNb, body, string(argsJSON))
}

// }}}
// {{{ Private functions

// }}}
// {{{ Package functions

// }}}
