package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.12.0"
)

var tp *trace.TracerProvider

type resourceConfig struct {
	serviceName    string
	serviceVersion string
	environment    string

	attributes map[string]string
}

// NewResource returns a resource describing this application.
func NewResource(opts ...Option) *resource.Resource {
	// 默认值
	rc := &resourceConfig{
		serviceName:    "demo-service",
		serviceVersion: "v0.0.0",
		environment:    "dev",
	}
	apply(rc, opts...)

	kvs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(rc.serviceName),
		semconv.ServiceVersionKey.String(rc.serviceVersion),
		attribute.String("env", rc.environment),
	}
	for k, v := range rc.attributes {
		kvs = append(kvs, attribute.String(k, v))
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, kvs...),
	)
	if err != nil {
		panic(err)
	}

	return r
}

// Init 初始化链路跟踪，fraction为分数，默认为1.0，值>=1.0表示全部链路都采样, 值<=0表示全部都不采样，0<值<1只采样百分比
func Init(exporter trace.SpanExporter, resource *resource.Resource, fractions ...float64) {
	var fraction = 1.0
	if len(fractions) > 0 {
		if fractions[0] <= 0 {
			fraction = 0
		} else if fractions[0] < 1 {
			fraction = fractions[0]
		}
	}

	tp = trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(fraction))), // 采样率
	)
	// 将TracerProvider注册为全局，这样将来任何导入包go.opentelemetry.io/otel/trace后，就可以默认使用它。
	otel.SetTracerProvider(tp)
}

// Close 停止
func Close(ctx context.Context) error {
	if tp == nil {
		return nil
	}
	return tp.Shutdown(ctx)
}
