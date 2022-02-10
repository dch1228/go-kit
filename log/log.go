package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	base *zap.Logger
}

func (lg *Logger) Debug(msg string, fields ...Field) {
	lg.base.Debug(msg, fields...)
}

func (lg *Logger) Info(msg string, fields ...Field) {
	lg.base.Info(msg, fields...)
}

func (lg *Logger) Warn(msg string, fields ...Field) {
	lg.base.Warn(msg, fields...)
}

func (lg *Logger) Error(msg string, fields ...Field) {
	lg.base.Error(msg, fields...)
}

func (lg *Logger) Sugar() *zap.SugaredLogger {
	return lg.base.WithOptions(AddCallerSkip(-1)).Sugar()
}

func (lg *Logger) Named(s string) *Logger {
	return &Logger{base: lg.base.Named(s)}
}

func (lg *Logger) With(fields ...Field) *Logger {
	return &Logger{base: lg.base.With(fields...)}
}

func (lg *Logger) WithOptions(opts ...Option) *Logger {
	return &Logger{base: lg.base.WithOptions(opts...)}
}

func (lg *Logger) Sync() error {
	return lg.base.Sync()
}

func New(cfg Config) (*Logger, error) {
	var (
		output    zapcore.WriteSyncer
		errOutput zapcore.WriteSyncer
	)
	if cfg.File.Filename != "" {
		lg := newFileLog(cfg.File)
		output = zapcore.AddSync(lg)
	} else {
		stdout, _, err := zap.Open([]string{"stdout"}...)
		if err != nil {
			return nil, err
		}
		output = stdout
	}

	if cfg.ErrorOutputPath != "" {
		errOut, _, err := zap.Open([]string{cfg.ErrorOutputPath}...)
		if err != nil {
			return nil, err
		}
		errOutput = errOut
	} else {
		errOutput = output
	}

	return newLoggerWithWriteSyncer(cfg, output, errOutput)
}

func newLoggerWithWriteSyncer(cfg Config, output, errOutput zapcore.WriteSyncer) (*Logger, error) {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, err
	}
	encoder := NewEncoder(cfg)
	if err != nil {
		return nil, err
	}
	core := zapcore.NewCore(encoder, output, level)

	opts := append([]Option{
		zap.ErrorOutput(errOutput),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	})

	lg := zap.New(core, opts...)
	return &Logger{base: lg}, nil
}

func NewEncoder(cfg Config) zapcore.Encoder {
	cc := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch cfg.Format {
	case "text", "":
		return zapcore.NewConsoleEncoder(cc)
	case "json":
		return zapcore.NewJSONEncoder(cc)
	default:
		panic(fmt.Sprintf("unsupport log format: %s", cfg.Format))
	}
}

func newFileLog(cfg FileConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  true,
	}
}
