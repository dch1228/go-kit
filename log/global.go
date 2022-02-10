package log

import (
	"sync"
)

var (
	mu     sync.Mutex
	global = MustNew(Config{
		Level: "info",
	})
)

func SetLogger(lg *Logger) {
	mu.Lock()
	defer mu.Unlock()

	global = lg
}

func Debug(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(1)).Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(1)).Sugar().Debugf(format, args...)
}

func Info(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(1)).Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(1)).Sugar().Infof(format, args...)
}

func Warn(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(1)).Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(1)).Sugar().Warnf(format, args...)
}

func Error(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(1)).Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(1)).Sugar().Errorf(format, args...)
}

func With(fields ...Field) *Logger {
	return global.With(fields...)
}

func L() *Logger {
	return global
}

func Sync() error {
	return global.Sync()
}
