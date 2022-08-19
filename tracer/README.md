## tracer

在[go.opentelemetry.io/otel](go.opentelemetry.io/otel)基础上封装的链路跟踪库。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/tracer

<br>

## 使用示例

```go
func initTrace() {
	exporter, err := tracer.NewJaegerExporter("http://localhost:14268/api/traces")
	if err != nil {
		panic(err)
	}

	resource := tracer.NewResource(
		tracer.WithServiceName("your-service-name"),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // 默认采集全部
	// tracer.Init(exporter, resource, 0.5) // 采集一半
}
```

<br>

documents https://opentelemetry.io/docs/instrumentation/go/

support OpenTelemetry in other libraries https://opentelemetry.io/registry/?language=go&component=instrumentation
