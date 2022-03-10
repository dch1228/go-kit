package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

type CtxField func(ctx context.Context) Field

func buildCtxField(ctx context.Context, fields ...CtxField) []Field {
	out := make([]Field, 0, len(fields))
	for _, field := range fields {
		out = append(out, field(ctx))
	}
	return out
}

func TraceID() CtxField {
	return func(ctx context.Context) Field {
		if span := trace.SpanContextFromContext(ctx); span.HasTraceID() {
			return String("trace_id", span.TraceID().String())
		}
		return Skip()
	}
}

func SpanID() CtxField {
	return func(ctx context.Context) Field {
		if span := trace.SpanContextFromContext(ctx); span.HasSpanID() {
			return String("span_id", span.SpanID().String())
		}
		return Skip()
	}
}
