package mysql

import "time"

var (
	defaultIsLog = false // 输出日志

	defaultMaxIdleConns    = 3               // 空闲连接数
	defaultMaxOpenConns    = 30              // 最大连接数
	defaultConnMaxLifetime = 5 * time.Minute // 5分钟后断开多余的空闲连接

	defaultTLSKey         = "custom" // 连接地址字段tls值
	defaultCAFile         = ""       // ca证书
	defaultClientKeyFile  = ""       // 客户端key
	defaultClientCertFile = ""       // 客户端cert

	defaultTables []interface{}
)

type options struct {
	IsLog           bool
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	tlsKey          string
	caFile          string
	clientKeyFile   string
	clientCertFile  string
	tables          []interface{}
}

type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func defaultOptions() *options {
	return &options{
		IsLog:           defaultIsLog,
		maxIdleConns:    defaultMaxIdleConns,
		maxOpenConns:    defaultMaxOpenConns,
		connMaxLifetime: defaultConnMaxLifetime,
		tlsKey:          defaultTLSKey,
		caFile:          defaultCAFile,
		clientKeyFile:   defaultClientKeyFile,
		clientCertFile:  defaultClientCertFile,
		tables:          nil,
	}
}

// WithLog set log sql
func WithLog() Option {
	return func(o *options) {
		o.IsLog = true
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

// WithTLSKey set tls key
func WithTLSKey(key string) Option {
	return func(o *options) {
		o.tlsKey = key
	}
}

// WithCAFile set ca file
func WithCAFile(file string) Option {
	return func(o *options) {
		o.caFile = file
	}
}

// WithClientKeyFile set client key file
func WithClientKeyFile(file string) Option {
	return func(o *options) {
		o.clientKeyFile = file
	}
}

// WithClientCertFile set client cert file
func WithClientCertFile(file string) Option {
	return func(o *options) {
		o.clientCertFile = file
	}
}

// WithTable set Auto Migrate tables
func WithTable(tables ...interface{}) Option {
	return func(o *options) {
		o.tables = tables
	}
}
