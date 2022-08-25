package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	// Print body max length
	defaultMaxLength = 300

	// default zap log
	defaultLogger, _ = zap.NewDevelopment() // 默认日志输出

	// Default request id name
	defaultRequestIDName = "X-Request-Id"

	// Ignore route list
	defaultIgnoreRoutes = map[string]struct{}{
		"/ping": {},
		"/pong": {},
	}
)

func defaultOptions() *options {
	return &options{
		maxLength:     defaultMaxLength,
		log:           defaultLogger,
		ignoreRoutes:  defaultIgnoreRoutes,
		requestIDName: "",
	}
}

type options struct {
	maxLength     int
	log           *zap.Logger
	ignoreRoutes  map[string]struct{}
	requestIDName string
}

// Option set the gin logger options.
type Option func(*options)

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMaxLen logger content max length
func WithMaxLen(maxLen int) Option {
	return func(o *options) {
		o.maxLength = maxLen
	}
}

// WithLog set log
func WithLog(log *zap.Logger) Option {
	return func(o *options) {
		if log != nil {
			o.log = log
		}
	}
}

// WithIgnoreRoutes no logger content routes
func WithIgnoreRoutes(routes ...string) Option {
	return func(o *options) {
		for _, route := range routes {
			o.ignoreRoutes[route] = struct{}{}
		}
	}
}

// WithRequestID name is field in header, eg:X-Request-Id
func WithRequestID(name ...string) Option {
	var requestIDName string
	if len(name) > 0 && name[0] != "" {
		requestIDName = name[0]
	} else {
		requestIDName = defaultRequestIDName
	}
	return func(o *options) {
		o.requestIDName = requestIDName
	}
}

// ------------------------------------------------------------------------------------------

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func getBodyData(buf *bytes.Buffer, maxLen int) string {
	var body string

	if buf.Len() > maxLen {
		body = string(buf.Bytes()[:maxLen]) + " ...... "
		// 如果有敏感数据需要过滤掉，比如明文密码
	} else {
		body = buf.String()
	}

	return body
}

// Logging print request and response info
func Logging(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		start := time.Now()

		// 忽略打印指定的路由
		if _, ok := o.ignoreRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		//  处理前打印输入信息
		buf := bytes.Buffer{}
		_, _ = buf.ReadFrom(c.Request.Body)

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
		}
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch || c.Request.Method == http.MethodDelete {
			fields = append(fields,
				zap.Int("size", buf.Len()),
				zap.String("body", getBodyData(&buf, o.maxLength)),
			)
		}
		reqID := ""
		if o.requestIDName != "" {
			reqID = c.Request.Header.Get(o.requestIDName)
			fields = append(fields, zap.String(o.requestIDName, reqID))
		}
		o.log.Info("<<<<", fields...)

		c.Request.Body = io.NopCloser(&buf)

		//  替换writer
		newWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = newWriter

		//  处理请求
		c.Next()

		// 处理后打印返回信息
		fields = []zap.Field{
			zap.Int("code", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.Path),
			zap.Int64("time_us", time.Since(start).Nanoseconds()/1000),
			zap.Int("size", newWriter.body.Len()),
			zap.String("response", strings.TrimRight(getBodyData(newWriter.body, o.maxLength), "\n")),
		}
		if o.requestIDName != "" {
			fields = append(fields, zap.String(o.requestIDName, reqID))
		}
		o.log.Info(">>>>", fields...)
	}
}
