package middleware

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"github.com/reugn/equalizer"
	"google.golang.org/grpc"
)

// ---------------------------------- server interceptor ----------------------------------

// RateLimitOption 日志设置
type RateLimitOption func(*rateLimitOptions)

type rateLimitOptions struct {
	capacity       int32         // 允许请求最大峰值
	qps            int64         // 允许请求速度
	refillInterval time.Duration // 填充token速度，refillInterval=time.Second/qps
}

func defaultRateLimitOptions() *rateLimitOptions {
	return &rateLimitOptions{
		capacity:       1000,
		qps:            500,
		refillInterval: time.Second / 500,
	}
}

func (o *rateLimitOptions) apply(opts ...RateLimitOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithRateLimitCapacity 设置允许请求最大峰值
func WithRateLimitCapacity(capacity int) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.capacity = int32(capacity)
	}
}

// WithRateLimitQPS 设置请求qps
func WithRateLimitQPS(qps int64) RateLimitOption {
	return func(o *rateLimitOptions) {
		o.qps = qps
		o.refillInterval = time.Second / time.Duration(o.qps)
	}
}

type myLimiter struct {
	TB *equalizer.TokenBucket // 令牌桶
}

func (m *myLimiter) Limit() bool {
	if m.TB.Ask() {
		return false
	}

	return true
}

// UnaryServerRateLimit 限流unary拦截器
func UnaryServerRateLimit(opts ...RateLimitOption) grpc.UnaryServerInterceptor {
	o := defaultRateLimitOptions()
	o.apply(opts...)

	limiter := &myLimiter{equalizer.NewTokenBucket(o.capacity, o.refillInterval)}
	return ratelimit.UnaryServerInterceptor(limiter)
}

// StreamServerRateLimit 限流stream拦截器
func StreamServerRateLimit(opts ...RateLimitOption) grpc.StreamServerInterceptor {
	o := defaultRateLimitOptions()
	o.apply(opts...)

	limiter := &myLimiter{equalizer.NewTokenBucket(o.capacity, o.refillInterval)}
	return ratelimit.StreamServerInterceptor(limiter)
}
