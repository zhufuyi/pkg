package ratelimiter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	// default qps value
	defaultQPS rate.Limit = 500

	// default the maximum instantaneous request spike allowed, burst >= qps
	defaultBurst = 1000

	// default is ip limit
	defaultIsIP = false
)

func defaultOptions() *options {
	return &options{
		qps:   defaultQPS,
		burst: defaultBurst,
		isIP:  false,
	}
}

type options struct {
	qps   rate.Limit
	burst int
	isIP  bool // false: path limit, true: IP limit
}

// Option logger middleware options
type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithQPS set the qps value
func WithQPS(qps rate.Limit) Option {
	return func(o *options) {
		o.qps = qps
	}
}

// WithBurst set the burst value, burst >= qps
func WithBurst(burst int) Option {
	return func(o *options) {
		o.burst = burst
	}
}

// WithPath set the path limit mode
func WithPath() Option {
	return func(o *options) {
		o.isIP = false
	}
}

// WithIP set the path limit mode
func WithIP() Option {
	return func(o *options) {
		o.isIP = true
	}
}

// QPS set limit qps parameters
func QPS(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)
	l = NewLimiter()

	return func(c *gin.Context) {
		var path string
		if !o.isIP {
			path = c.FullPath()
		} else {
			path = c.ClientIP()
		}

		l.qpsLimiter.LoadOrStore(path, rate.NewLimiter(o.qps, o.burst))
		if !l.allow(path) {
			c.JSON(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))
			return
		}
		c.Next()
	}
}
