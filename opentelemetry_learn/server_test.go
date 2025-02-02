package opentelemetrylearn

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func TestServer(t *testing.T) {
	res, err := newResource("demo", "v0.0.1")
	require.NoError(t, err)

	prop := newPropagator()
	// passing tracing info from client and server
	otel.SetTextMapPropagator(prop)

	tp, err := newTraceProvider(res)
	require.NoError(t, err)
	defer tp.Shutdown(context.Background())
	otel.SetTracerProvider(tp)

	server := gin.Default()
	server.GET("/", func(ginCtx *gin.Context) {
		tracer := otel.Tracer("github.com/blabla/opentelemetry")
		// var ctx context.Context = ginCtx
		ctx, span := tracer.Start(ginCtx, "top-span")
		defer span.End()

		time.Sleep(time.Second)

		span.AddEvent("something happened")

		ctx, subSpan := tracer.Start(ctx, "sub-span")
		defer subSpan.End()

		subSpan.SetAttributes(attribute.String("attr1", "value1"))

		time.Sleep(time.Millisecond * 300)

		ginCtx.String(http.StatusOK, "test")
	})
	server.Run(":8081")
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(res *resource.Resource) (*trace.TracerProvider, error) {
	exporter, err := zipkin.New(
		"http://localhost:9411/api/v2/spans",
	)
	if err != nil {
		return nil, err
	}
	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter, trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}
