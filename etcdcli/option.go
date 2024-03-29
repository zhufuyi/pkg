package etcdcli

import (
	"time"

	"go.uber.org/zap"
)

// Option set the etcd client options.
type Option func(*options)

type options struct {
	dialTimeout time.Duration // 连接超时时间，单位(秒)

	username string
	password string

	isSecure           bool
	serverNameOverride string // etcd 域名
	certFile           string // 证书文件路径

	autoSyncInterval time.Duration // 自动同步member list的间隔
	logger           *zap.Logger
}

func defaultOptions() *options {
	return &options{
		dialTimeout: time.Second * 5,
	}
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithDialTimeout set dial timeout
func WithDialTimeout(duration time.Duration) Option {
	return func(o *options) {
		o.dialTimeout = duration
	}
}

// WithAuth set authentication
func WithAuth(username string, password string) Option {
	return func(o *options) {
		o.username = username
		o.password = password
	}
}

// WithSecure set tls
func WithSecure(serverNameOverride string, certFile string) Option {
	return func(o *options) {
		o.isSecure = true
		o.serverNameOverride = serverNameOverride
		o.certFile = certFile
	}
}

// WithAutoSyncInterval set auto sync interval value
func WithAutoSyncInterval(duration time.Duration) Option {
	return func(o *options) {
		o.autoSyncInterval = duration
	}
}

// WithLog set logger
func WithLog(l *zap.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}
