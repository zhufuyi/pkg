package mysql

import (
	"time"
)

var (
	defaultIsLog                       = false // 是否输出日志
	defaultSlowThreshold time.Duration = 0     // 如果大于0，只打印时间大于阈值的日志，优先级比isLog高

	defaultMaxIdleConns    = 3               // 空闲连接数
	defaultMaxOpenConns    = 30              // 最大连接数
	defaultConnMaxLifetime = 5 * time.Minute // 5分钟后断开多余的空闲连接
)

type options struct {
	isLog           bool
	slowThreshold   time.Duration
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
}

type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultOptions() *options {
	return &options{
		isLog:           defaultIsLog,
		slowThreshold:   defaultSlowThreshold,
		maxIdleConns:    defaultMaxIdleConns,
		maxOpenConns:    defaultMaxOpenConns,
		connMaxLifetime: defaultConnMaxLifetime,
	}
}

// WithLog set log sql
func WithLog() Option {
	return func(o *options) {
		o.isLog = true
	}
}

// WithSlowThreshold Set sql values greater than the threshold
func WithSlowThreshold(d time.Duration) Option {
	return func(o *options) {
		o.slowThreshold = d
	}
}

// WithMaxIdleConns set max idle conns
func WithMaxIdleConns(size int) Option {
	return func(o *options) {
		o.maxIdleConns = size
	}
}

// WithMaxOpenConns set max open conns
func WithMaxOpenConns(size int) Option {
	return func(o *options) {
		o.maxOpenConns = size
	}
}

// WithConnMaxLifetime set conn max liftime
func WithConnMaxLifetime(t time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = t
	}
}
