package prometheus

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceID string
	Help       string
}

func NewPrometheusBuilder(namespace, subSystem, name, instanceID, help string) *Builder {
	return &Builder{
		Namespace:  namespace,
		Subsystem:  subSystem,
		Name:       name,
		InstanceID: instanceID,
		Help:       help,
	}
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	// pattern: route
	labels := []string{"method", "pattern", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_resp_time",
		ConstLabels: prometheus.Labels{
			"instance_id": b.InstanceID,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.75:  0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
		Help: b.Help,
	}, labels)
	if err := prometheus.Register(vector); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			vector = are.ExistingCollector.(*prometheus.SummaryVec)
		}
	}
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration))
		}()
		ctx.Next()
	}
}

func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Name:      b.Name + "_active_req",
		Help:      b.Help,
		ConstLabels: prometheus.Labels{
			"instance_id": b.InstanceID,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}
}
