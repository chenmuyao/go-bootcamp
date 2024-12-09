package ratelimit

import (
	"context"
	"errors"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/chenmuyao/go-bootcamp/pkg/limiter"
)

var errLimited = errors.New("sms queries reach the limit")

type RateLimitSMSService struct {
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

func NewRateLimitSMSService(svc sms.Service, l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}

func (r *RateLimitSMSService) Send(
	ctx context.Context,
	toNb string,
	body string,
	args ...string,
) error {
	accepted := r.limiter.AcceptConnection(ctx, r.key)
	if !accepted {
		// limited
		return errLimited
	}
	return r.svc.Send(ctx, toNb, body, args...)
}
