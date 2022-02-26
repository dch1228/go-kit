package log

import (
	"context"
	"fmt"
)

type Logger struct {
	base      *zapLogger
	ctx       context.Context
	ctxFields []CtxField
}

func New(cfg Config) *Logger {
	var (
		output    WriteSyncer
		errOutput WriteSyncer
	)
	if cfg.File.Filename != "" {
		lg := newFileLog(cfg.File)
		output = AddSync(lg)
	} else {
		stdout, _, err := Open([]string{"stdout"}...)
		if err != nil {
			panic(fmt.Sprintf("log init error: %s", err))
		}
		output = stdout
	}

	if cfg.ErrorOutputPath != "" {
		errOut, _, err := Open([]string{cfg.ErrorOutputPath}...)
		if err != nil {
			panic(fmt.Sprintf("log init error: %s", err))
		}
		errOutput = errOut
	} else {
		errOutput = output
	}

	level := NewAtomicLevel()
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		panic(fmt.Sprintf("log level error: %s", err))
	}

	opts := append([]zapOption{
		ErrorOutput(errOutput),
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
	})

	return &Logger{
		base: zapNew(
			NewCore(newEncoder(cfg), output, level),
			opts...,
		),
		ctx: context.Background(),
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
