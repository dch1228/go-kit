package log

import (
	"context"
	"testing"
)

func TestCtxField(t *testing.T) {
	lg := New(Config{
		Level: "info",
	})

	lg = lg.WithCtxFields(
		func() CtxField {
			return func(ctx context.Context) Field {
				val := ctx.Value("test")
				if val != nil {
					return String("ctxfield", val.(string))
				}
				return Skip()
			}
		}())

	ctx := context.WithValue(context.Background(), "test", "test")
	lg.WithCtx(ctx).Info("test")
	lg.WithCtx(context.Background()).Info("test")
}
