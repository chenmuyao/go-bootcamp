package opentelemetry

import (
	"context"

	"github.com/chenmuyao/go-bootcamp/internal/service/sms"
	"go.opentelemetry.io/otel/trace"
)

type OTELSMSSvc struct {
	svc    sms.Service
	tracer trace.Tracer
}

// Send implements sms.Service.
func (o *OTELSMSSvc) Send(ctx context.Context, toNb string, body string, args ...string) error {
	ctx, span := o.tracer.Start(ctx, "sms")
	defer span.End()
	span.AddEvent("send sms")
	err := o.svc.Send(ctx, toNb, body, args...)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

func NewOTELSMSSvc(svc sms.Service, tracer trace.Tracer) sms.Service {
	return &OTELSMSSvc{
		svc:    svc,
		tracer: tracer,
	}
}
