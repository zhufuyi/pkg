package tracer

import (
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewJaegerExporter use jaeger collector as exporter, e.g. default url=http://localhost:14268/api/traces
func NewJaegerExporter(url string) (sdkTrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(url),
			//jaeger.WithUsername("your-username"),
			//jaeger.WithPassword("your-password"),
		),
	)
}

// NewJaegerAgentExporter use jaeger agent as exporter
func NewJaegerAgentExporter(host string, port string) (sdkTrace.SpanExporter, error) {
	return jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host),
			jaeger.WithAgentPort(port),
		),
	)
}
