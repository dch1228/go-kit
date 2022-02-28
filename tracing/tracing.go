package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/dch1228/go-kit/log"
)

func Init(cfg Config) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Endpoint)))
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(cfg.SamplerRatio))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(cfg.Name),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func TraceID() log.CtxField {
	return func(ctx context.Context) log.Field {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return log.String("trace_id", span.TraceID().String())
		}
		return log.Skip()
	}
}

func SpanID() log.CtxField {
	return func(ctx context.Context) log.Field {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return log.String("span_id", span.SpanID().String())
		}
		return log.Skip()
	}
}
