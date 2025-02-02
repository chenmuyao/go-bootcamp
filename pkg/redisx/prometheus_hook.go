package redisx

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type PrometheusHook struct {
	vector *prometheus.SummaryVec
}

func NewPrometheusHook(opts prometheus.SummaryOpts) redis.Hook {
	vector := prometheus.NewSummaryVec(opts, []string{
		"cmd", "key_exist",
	})
	if err := prometheus.Register(vector); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			vector = are.ExistingCollector.(*prometheus.SummaryVec)
		}
	}
	return &PrometheusHook{}
}

// DialHook implements redis.Hook.
func (p *PrometheusHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook implements redis.Hook.
func (p *PrometheusHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		var err error
		start := time.Now()
		defer func() {
			// NOTE: add biz
			// biz := ctx.Value("biz")
			duration := time.Since(start).Milliseconds()
			keyExists := err == redis.Nil
			p.vector.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExists)).
				Observe(float64(duration))
		}()
		err = next(ctx, cmd)
		return err
	}
}

// ProcessPipelineHook implements redis.Hook.
func (p *PrometheusHook) ProcessPipelineHook(
	next redis.ProcessPipelineHook,
) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
