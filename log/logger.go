package log

import (
	"context"
	"fmt"
	"sync"

	klog "github.com/go-kratos/kratos/v2/log"
)

var (
	once   = sync.Once{}
	global *Logger
)

type Logger struct {
	base      *zapLogger
	ctx       context.Context
	ctxFields []CtxField
}

func Setup(cfg Config) (err error) {
	once.Do(func() {
		err = setup(cfg)
	})
	return
}

func setup(cfg Config) error {
	var (
		output WriteSyncer
	)

	if cfg.File.Filename != "" {
		lg := newFileLog(cfg.File)
		output = AddSync(lg)
	} else {
		stdout, _, err := Open([]string{"stdout"}...)
		if err != nil {
			return err
		}
		output = stdout
	}

	level := NewAtomicLevel()
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}

	opts := append([]zapOption{
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
	})

	var ctxFields []CtxField
	if cfg.EnableTrace {
		ctxFields = append(ctxFields, TraceID(), SpanID())
	}

	logger := &Logger{
		base: zapNew(
			NewCore(newEncoder(cfg), output, level),
			opts...,
		),
		ctx:       context.Background(),
		ctxFields: ctxFields,
	}

	global = logger
	klog.SetLogger(logger)
	return nil
}

func MustSetup(cfg Config) {
	if err := Setup(cfg); err != nil {
		panic(err)
	}
}

func newEncoder(cfg Config) Encoder {
	cc := EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     DefaultLineEnding,
		EncodeLevel:    CapitalLevelEncoder,
		EncodeTime:     RFC3339TimeEncoder,
		EncodeDuration: MillisDurationEncoder,
		EncodeCaller:   ShortCallerEncoder,
	}
	switch cfg.Format {
	case "text", "":
		return NewConsoleEncoder(cc)
	case "json":
		return NewJSONEncoder(cc)
	default:
		panic("create encoder error: unsupported format")
	}
}

func newFileLog(cfg FileConfig) *fileLogger {
	return &fileLogger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  true,
	}
}

func (log *Logger) Named(s string) *Logger {
	return &Logger{
		base:      log.base.Named(s),
		ctx:       log.ctx,
		ctxFields: log.ctxFields,
	}
}

func (log *Logger) Info(msg string, fields ...Field) {
	fields = append(fields, buildCtxField(log.ctx, log.ctxFields...)...)
	log.base.Info(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...Field) {
	fields = append(fields, buildCtxField(log.ctx, log.ctxFields...)...)
	log.base.Warn(msg, fields...)
}

func (log *Logger) Error(msg string, err error, fields ...Field) {
	fields = append(fields, buildCtxField(log.ctx, log.ctxFields...)...)
	fields = append(fields, Err(err))
	log.base.Error(msg, fields...)
}

func (log *Logger) With(fields ...Field) *Logger {
	return &Logger{
		base:      log.base.With(fields...),
		ctx:       log.ctx,
		ctxFields: log.ctxFields,
	}
}

func (log *Logger) WithOptions(options ...zapOption) *Logger {
	return &Logger{
		base:      log.base.WithOptions(options...),
		ctx:       log.ctx,
		ctxFields: log.ctxFields,
	}
}

func (log *Logger) WithCtx(ctx context.Context) *Logger {
	return &Logger{
		base:      log.base,
		ctx:       ctx,
		ctxFields: log.ctxFields,
	}
}

func (log *Logger) WithCtxFields(fields ...CtxField) *Logger {
	return &Logger{
		base:      log.base,
		ctx:       log.ctx,
		ctxFields: append(log.ctxFields, fields...),
	}
}

func (log *Logger) Sync() error {
	return log.base.Sync()
}

// Log 实现 kratos logger
func (log *Logger) Log(level klog.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		log.Warn("Key values must appear in pairs", Any("keyvals", keyvals))
		return nil
	}

	var data []Field
	for i := 0; i < len(keyvals); i += 2 {
		data = append(data, Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case klog.LevelDebug, klog.LevelInfo:
		log.Info("", data...)
	case klog.LevelWarn:
		log.Warn("", data...)
	case klog.LevelError, klog.LevelFatal:
		log.Error("", nil, data...)
	}
	return nil
}
