package render

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/zhufuyi/pkg/logger"
	"go.uber.org/zap"

	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

var (
	tracer opentracing.Tracer
)

func initTrace(serviceName string, jaegerHostPort string) {
	cfg := &jaegerConfig.Configuration{
		ServiceName: serviceName,

		Sampler: &jaegerConfig.SamplerConfig{
			Type:  "const", // 固定采样
			Param: 1,       // 1:全采样，0:不采样
		},

		Reporter: &jaegerConfig.ReporterConfig{
			LogSpans:           false, // 是否在终端输出Reporting span信息
			LocalAgentHostPort: jaegerHostPort,
		},
	}

	jaegerLogger := jaegerLoggerAdapter{log: logger.WithFields(logger.String("service", serviceName))}

	var err error
	tracer, _, err = cfg.NewTracer(jaegerConfig.Logger(jaegerLogger))
	if err != nil {
		tracer = nil
		logger.Error("init jaegertrace failed", logger.Err(err), logger.String("jaegerHostPort", jaegerHostPort))
		return
	}

	opentracing.SetGlobalTracer(tracer)

	logger.Info("init jaeger trace success.")
}

func HttpTrace(IsEnable bool, appName string, addr string) gin.HandlerFunc {
	var parentSpan opentracing.Span

	return func(c *gin.Context) {
		if IsEnable {
			// 判断是否已经初始化过tracer
			if tracer == nil {
				initTrace(appName, addr)
			}

			spCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
			if err != nil {
				parentSpan = tracer.StartSpan(c.Request.URL.Path)
				parentSpan.SetTag(string(ext.Component), "HTTP")
				parentSpan.SetTag(string(ext.HTTPMethod), c.Request.Method)
				parentSpan.SetTag(string(ext.HTTPUrl), c.Request.RequestURI)
			} else {
				parentSpan = opentracing.StartSpan(
					c.Request.URL.Path,
					opentracing.ChildOf(spCtx),
					opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
					opentracing.Tag{Key: string(ext.HTTPMethod), Value: c.Request.Method},
					opentracing.Tag{Key: string(ext.HTTPUrl), Value: c.Request.RequestURI},
					ext.SpanKindRPCServer,
				)
			}
			defer parentSpan.Finish()

			// 通过gin.Context的字段Keys传递跟踪信息给下一个span
			c.Keys = map[string]interface{}{
				"tracer":        tracer,
				"parentSpan":    parentSpan,
				"isEnableTrace": IsEnable,
			}
		}

		c.Next()

		if !IsEnable || parentSpan == nil {
			return
		}

		code, ok := c.Get("code")
		if ok && code != nil {
			httpStatus := code.(int)
			parentSpan.SetTag(string(ext.HTTPStatusCode), httpStatus) // 高并发请求同一个接口，多个span同样的tag会可能会出现警告
			if httpStatus != http.StatusOK {
				if err, ok := c.Get("error"); ok {
					ext.Error.Set(parentSpan, true)
					//parentSpan.LogKV("event", fmt.Sprintf("%s=%d error=%+v", ext.HTTPStatusCode, httpStatus, err))
					parentSpan.LogKV("event", "http request error", "error", err)
				}
			}
		}
	}
}

type jaegerLoggerAdapter struct {
	log *zap.Logger
}

func (l jaegerLoggerAdapter) Error(msg string) {
	l.log.Error(msg)
}

func (l jaegerLoggerAdapter) Infof(msg string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(msg, args...))
}
