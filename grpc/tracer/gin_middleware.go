package tracer

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/zhufuyi/pkg/grpc/tracer/otgrpc"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// GinCtx 把gin的context转换为标准的context
func GinCtx(c *gin.Context) context.Context {
	tracerVal, _ := c.Get(otgrpc.GinTracerKey)
	ctx := context.WithValue(context.Background(), otgrpc.GinTracerKey, tracerVal) //nolint

	parentSpanVal, _ := c.Get(otgrpc.GinParentSpanKey)
	return context.WithValue(ctx, otgrpc.GinParentSpanKey, parentSpanVal) //nolint
}

// GinMiddleware gin的链路跟踪中间件
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startSpan := opentracing.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		// 通过gin的context传递的tracer和startSpan给下一个使用者
		c.Set(otgrpc.GinTracerKey, opentracing.GlobalTracer())
		c.Set(otgrpc.GinParentSpanKey, startSpan)

		c.Next()
	}
}

// Get 在http请求中添加链路追踪
func Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	span, newCtx := opentracing.StartSpanFromContext(ctx, "HTTP GET: "+url, opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
	defer span.Finish()

	req = req.WithContext(newCtx)
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() //nolint

	return io.ReadAll(resp.Body)
}
