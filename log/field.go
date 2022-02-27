package log

import (
	"context"
)

type CtxField func(ctx context.Context) Field

func buildCtxField(ctx context.Context, fields ...CtxField) []Field {
	out := make([]Field, 0, len(fields))
	for _, field := range fields {
		out = append(out, field(ctx))
	}
	return out
}
