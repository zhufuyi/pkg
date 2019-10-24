package render

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

// 忽略部分路由打印
var ignoreRoutes map[string]bool

func init() {
	ignoreRoutes = map[string]bool{
		"/getSubID": true,
	}
}

// 限制显示body内容最大长度
const limitSize = 300

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// InOutLog gin输入输出日志
func InOutLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 忽略打印指定的路由
		if isIgnoreRoute(c.Request.URL.Path) {
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
				logger.String("body", getBodyData(&buf)),
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
			logger.String("response", strings.TrimRight(getBodyData(newWriter.body), "\n")),
		)
	}
}

func getBodyData(buf *bytes.Buffer) string {
	var body string

	if buf.Len() > limitSize {
		body = string(buf.Bytes()[:limitSize]) + " ...... "
		// 如果有敏感数据需要过滤掉，比如明文密码
	} else {
		body = buf.String()
	}

	return body
}

func isIgnoreRoute(routeValue string) bool {
	for route, v := range ignoreRoutes {
		if strings.Contains(routeValue, route) {
			return v
		}
	}

	return false
}
