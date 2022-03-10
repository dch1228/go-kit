package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLineEnding = zapcore.DefaultLineEnding

	InfoLevel  = zap.InfoLevel
	WarnLevel  = zap.WarnLevel
	ErrorLevel = zap.ErrorLevel
)

type (
	Field         = zapcore.Field
	Level         = zapcore.Level
	WriteSyncer   = zapcore.WriteSyncer
	Encoder       = zapcore.Encoder
	EncoderConfig = zapcore.EncoderConfig

	zapOption = zap.Option
	zapLogger = zap.Logger

	fileLogger = lumberjack.Logger
)

var (
	Any        = zap.Any
	Skip       = zap.Skip
	ByteString = zap.ByteString
	Bool       = zap.Bool
	Duration   = zap.Duration
	Float64    = zap.Float64
	Int        = zap.Int
	Int32      = zap.Int32
	Int64      = zap.Int64
	String     = zap.String
	Time       = zap.Time
	Uint       = zap.Uint
	Err        = zap.Error

	AddSync               = zapcore.AddSync
	NewCore               = zapcore.NewCore
	CapitalLevelEncoder   = zapcore.CapitalLevelEncoder
	RFC3339TimeEncoder    = zapcore.RFC3339TimeEncoder
	MillisDurationEncoder = zapcore.MillisDurationEncoder
	ShortCallerEncoder    = zapcore.ShortCallerEncoder
	NewConsoleEncoder     = zapcore.NewConsoleEncoder
	NewJSONEncoder        = zapcore.NewJSONEncoder

	Open           = zap.Open
	NewAtomicLevel = zap.NewAtomicLevel
	AddCaller      = zap.AddCaller
	AddCallerSkip  = zap.AddCallerSkip
	AddStacktrace  = zap.AddStacktrace
	WithCaller     = zap.WithCaller
	zapNew         = zap.New
)
