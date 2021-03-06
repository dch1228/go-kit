package log

import (
	"context"
)

func Info(msg string, fields ...Field) {
	global.base.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	global.base.Warn(msg, fields...)
}

func Error(msg string, err error, fields ...Field) {
	fields = append(fields, Err(err))
	global.base.Error(msg, fields...)
}

func With(fields ...Field) *Logger {
	return &Logger{
		base:      global.base.With(fields...),
		ctx:       global.ctx,
		ctxFields: global.ctxFields,
	}
}

func WithCtx(ctx context.Context) *Logger {
	return global.WithCtx(ctx)
}

func Sync() error {
	return global.Sync()
}

func L() *Logger {
	return global
}
