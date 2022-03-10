package log

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestCtxField(t *testing.T) {
	otel.SetTracerProvider(tracesdk.NewTracerProvider())

	MustSetup(Config{
		Level:       "info",
		EnableTrace: true,
	})

	tr := otel.Tracer("TestCtxField")
	ctx, span := tr.Start(context.Background(), "TestCtxField")
	defer span.End()

	WithCtx(ctx).Info("test")
}
