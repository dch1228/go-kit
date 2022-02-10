package log

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Field = zap.Field

func Skip() Field {
	return zap.Skip()
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func String(key string, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func Err(err error) Field {
	return zap.Error(err)
}

func TraceID(ctx context.Context) Field {
	return zap.String("trace_id", trace.SpanContextFromContext(ctx).TraceID().String())
}

func SpanID(ctx context.Context) Field {
	return zap.String("span_id", trace.SpanContextFromContext(ctx).SpanID().String())
}
