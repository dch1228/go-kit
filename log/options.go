package log

import "go.uber.org/zap"

type Option = zap.Option

func AddCallerSkip(skip int) Option {
	return zap.AddCallerSkip(skip)
}
