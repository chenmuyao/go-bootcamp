package prometheus

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusSMS struct {
	sms.Service
	sum prometheus.Summary
}

func (p *PrometheusSMS) Send(ctx context.Context, toNb string, body string, args ...string) error {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		p.sum.Observe(float64(duration))
	}()
	return p.Service.Send(ctx, toNb, body, args...)
}

func NewPrometheusOAuth(svc sms.Service, opts prometheus.SummaryOpts) sms.Service {
	sum := prometheus.NewSummary(opts)
	if err := prometheus.Register(sum); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			sum = are.ExistingCollector.(prometheus.Summary)
		}
	}
	return &PrometheusSMS{
		Service: svc,
		sum:     sum,
	}
}
