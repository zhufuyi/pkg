package middleware

import (
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ---------------------------------- client interceptor ----------------------------------

var (
	defaultTimes    uint = 3                                                                       // 重试次数
	defaultInterval      = time.Millisecond * 50                                                   // 重试间隔50毫秒
	defaultErrCodes      = []codes.Code{codes.Internal, codes.DeadlineExceeded, codes.Unavailable} // 默认触发重试的错误码
)

type options struct {
	times    uint
	interval time.Duration
	errCodes []codes.Code
}

func defaultOptions() *options {
	return &options{
		times:    defaultTimes,
		interval: defaultInterval,
		errCodes: defaultErrCodes,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the retry options.
type Option func(*options)

// WithRetryTimes 设置重试次数，最大10次
func WithRetryTimes(n uint) Option {
	return func(o *options) {
		if n > 10 {
			n = 10
		}
		o.times = n
	}
}

// WithRetryInterval 设置重试时间间隔，范围1毫秒到10秒
func WithRetryInterval(t time.Duration) Option {
	return func(o *options) {
		if t < time.Millisecond {
			t = time.Millisecond
		} else if t > 10*time.Second {
			t = 10 * time.Second
		}
		o.interval = t
	}
}

// WithRetryErrCodes 设置触发重试错误码
func WithRetryErrCodes(errCodes ...codes.Code) Option {
	for _, errCode := range errCodes {
		switch errCode {
		case codes.Internal, codes.DeadlineExceeded, codes.Unavailable:
		default:
			defaultErrCodes = append(defaultErrCodes, errCode)
		}
	}
	return func(o *options) {
		o.errCodes = defaultErrCodes
	}
}

// UnaryClientRetry 重试
func UnaryClientRetry(opts ...Option) grpc.UnaryClientInterceptor {
	o := defaultOptions()
	o.apply(opts...)
	return grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(o.times), // 设置重试次数
		grpc_retry.WithBackoff(func(attempt uint) time.Duration { // 设置重试间隔
			return o.interval
		}),
		grpc_retry.WithCodes(o.errCodes...), // 设置重试错误码
	)
}
