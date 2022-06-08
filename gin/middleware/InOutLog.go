package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/logger"
)

var (
	// 限制显示body内容最大长度
	defaultMaxLength = 300

	// 忽略打印路由
	defaultIgnoreRoutes = map[string]struct{}{
		"/ping": struct{}{},
		"/pong": struct{}{},
	}
)

func defaultOptions() *options {
	return &options{
		maxLength:    defaultMaxLength,
		ignoreRoutes: defaultIgnoreRoutes,
	}
}

type options struct {
	maxLength    int
	ignoreRoutes map[string]struct{}
}

// Option logger middleware options
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

// WithIgnoreRoutes no logger content routes
func WithIgnoreRoutes(routes ...string) Option {
	return func(o *options) {
		for _, route := range routes {
			o.ignoreRoutes[route] = struct{}{}
		}
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

// InOutLog gin输入输出日志
func InOutLog(opts ...Option) gin.HandlerFunc {
	o := defaultOptions()
	o.apply(opts...)

	return func(c *gin.Context) {
		// 忽略打印指定的路由
		if _, ok := o.ignoreRoutes[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		start := time.Now()

		//  处理前打印输入信息
		buf := bytes.Buffer{}
		buf.ReadFrom(c.Request.Body)
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			logger.Info("<<<<<<",
				logger.String("method", c.Request.Method),
				logger.Any("url", c.Request.URL),
				logger.Int("size", buf.Len()),
				logger.String("body", getBodyData(&buf, o.maxLength)),
			)
		} else {
			logger.Info("<<<<<<",
				logger.String("method", c.Request.Method),
				logger.Any("url", c.Request.URL),
			)
		}
		c.Request.Body = ioutil.NopCloser(&buf)

		//  替换writer
		newWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = newWriter

		//  处理请求
		c.Next()

		// 处理后打印返回信息
		logger.Info(">>>>>>",
			logger.Int("code", c.Writer.Status()),
			logger.String("method", c.Request.Method),
			logger.String("url", c.Request.URL.Path),
			logger.String("time", fmt.Sprintf("%dus", time.Now().Sub(start).Nanoseconds()/1000)),
			logger.Int("size", newWriter.body.Len()),
			logger.String("response", strings.TrimRight(getBodyData(newWriter.body, o.maxLength), "\n")),
		)
	}
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
