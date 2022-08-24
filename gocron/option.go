package gocron

import "go.uber.org/zap"

type options struct {
	zapLog *zap.Logger
}

func defaultOptions() *options {
	return &options{
		zapLog: nil,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the cron options.
type Option func(*options)

// WithLog 设置日志
func WithLog(log *zap.Logger) Option {
	return func(o *options) {
		o.zapLog = log
	}
}
