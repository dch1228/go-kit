package log

import (
	"fmt"
	"sync"
)

var (
	mu     sync.Mutex
	global *Logger
)

func init() {
	cfg := Config{
		Level: "info",
	}
	lg, err := New(cfg)
	if err != nil {
		panic(fmt.Sprintf("init log error: %s", err))
	}
	SetLogger(lg)
}

func SetLogger(lg *Logger) {
	mu.Lock()
	defer mu.Unlock()

	global = lg
}

func Debug(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(2)).Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(2)).Sugar().Debugf(format, args...)
}

func Info(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(2)).Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(2)).Sugar().Infof(format, args...)
}

func Warn(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(2)).Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(2)).Sugar().Warnf(format, args...)
}

func Error(msg string, fields ...Field) {
	global.WithOptions(AddCallerSkip(2)).Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	global.WithOptions(AddCallerSkip(2)).Sugar().Errorf(format, args...)
}

func With(fields ...Field) *Logger {
	return global.WithOptions(AddCallerSkip(1)).With(fields...)
}

func Sync() error {
	return global.Sync()
}
