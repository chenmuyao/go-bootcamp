package prometheus

import (
	"context"
	"time"

	"github.com/chenmuyao/go-bootcamp/internal/domain"
	"github.com/chenmuyao/go-bootcamp/internal/service/oauth2/gitea"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusOAuth struct {
	gitea.Service
	sum prometheus.Summary
}

func (p *PrometheusOAuth) VerifyCode(ctx context.Context, code string) (domain.GiteaInfo, error) {
	start := time.Now()
	defer func() {
		duration := time.Since(start).Milliseconds()
		p.sum.Observe(float64(duration))
	}()
	return p.Service.VerifyCode(ctx, code)
}

func NewPrometheusOAuth(svc gitea.Service, opts prometheus.SummaryOpts) gitea.Service {
	sum := prometheus.NewSummary(opts)
	if err := prometheus.Register(sum); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			sum = are.ExistingCollector.(prometheus.Summary)
		}
	}
	return &PrometheusOAuth{
		Service: svc,
		sum:     sum,
	}
}
