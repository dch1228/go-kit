package log

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported format")
)

type Logger = zap.Logger

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
			panic(fmt.Sprintf("log init error: %s", err))
		}
		output = stdout
	}

	if cfg.ErrorOutputPath != "" {
		errOut, _, err := zap.Open([]string{cfg.ErrorOutputPath}...)
		if err != nil {
			panic(fmt.Sprintf("log init error: %s", err))
		}
		errOutput = errOut
	} else {
		errOutput = output
	}

	return newLoggerWithWriteSyncer(cfg, output, errOutput)
}

func MustNew(cfg Config) *Logger {
	lg, err := New(cfg)
	if err != nil {
		panic(err)
	}
	return lg
}

func newLoggerWithWriteSyncer(cfg Config, output, errOutput zapcore.WriteSyncer) (*Logger, error) {
	level := zap.NewAtomicLevel()
	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return nil, err
	}
	encoder, err := newEncoder(cfg)
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
	return lg, nil
}

func newEncoder(cfg Config) (zapcore.Encoder, error) {
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
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch cfg.Format {
	case "text", "":
		return zapcore.NewConsoleEncoder(cc), nil
	case "json":
		return zapcore.NewJSONEncoder(cc), nil
	default:
		return nil, ErrUnsupportedFormat
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
