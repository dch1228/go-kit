package tracing

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	"github.com/dch1228/go-kit/log"
)

func TestLog(t *testing.T) {
	otel.SetTracerProvider(tracesdk.NewTracerProvider())

	lg := log.L().WithCtxFields(TraceID(), SpanID())

	tr := otel.Tracer("TestLog")
	ctx, span := tr.Start(context.Background(), "test")
	defer span.End()

	lg.WithCtx(ctx).Info("info")
}
