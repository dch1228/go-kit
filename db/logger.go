package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"

	"github.com/dch1228/go-kit/log"
)

var _ glog.Interface = (*Logger)(nil)

var ErrRecordNotFound = gorm.ErrRecordNotFound

type Logger struct {
	LogConfig
	base *log.Logger
}

func NewLogger(cfg LogConfig) *Logger {
	return &Logger{
		LogConfig: cfg,
		base:      log.L().Named("[GORM]"),
	}
}

func (l *Logger) LogMode(level glog.LogLevel) glog.Interface {
	newlogger := *l
	newlogger.Level = level
	return &newlogger
}

func (l *Logger) Info(_ context.Context, str string, args ...interface{}) {
	if l.Level < glog.Info {
		return
	}
	l.base.WithOptions(log.WithCaller(false)).Info(fmt.Sprintf(str, args...))
}

func (l *Logger) Warn(_ context.Context, str string, args ...interface{}) {
	if l.Level < glog.Warn {
		return
	}
	l.base.WithOptions(log.WithCaller(false)).Info(fmt.Sprintf(str, args...))
}

func (l *Logger) Error(_ context.Context, str string, args ...interface{}) {
	if l.Level < glog.Error {
		return
	}
	l.base.WithOptions(log.WithCaller(false)).Info(fmt.Sprintf(str, args...))
}

func (l *Logger) Trace(_ context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.Level <= glog.Silent {
		return
	}

	elapsed := time.Since(begin)
	lg := l.base.WithOptions(log.AddCallerSkip(2))
	switch {
	case err != nil && l.Level >= glog.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		lg.Error("trace", err, log.Duration("elapsed", elapsed), log.Int64("rows", rows), log.String("sql", sql))
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.Level >= glog.Warn:
		sql, rows := fc()
		lg.Warn("slow-sql", log.Duration("elapsed", elapsed), log.Int64("rows", rows), log.String("sql", sql))
	case l.Level == glog.Info:
		sql, rows := fc()
		lg.Info("trace", log.Duration("elapsed", elapsed), log.Int64("rows", rows), log.String("sql", sql))
	}
}
