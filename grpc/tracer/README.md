## tracer

支持rpc-->rpc链路跟踪，同时也支持gin-->rpc链路跟踪。

<br>

### 使用示例

#### grpc server

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 链路跟踪拦截器
	options = append(options, grpc.UnaryInterceptor(
		tracer.UnaryServerTracing(),
	))

	return options
}

func main() {
	// 连接jaeger服务端
	closer, _ := tracer.InitJaeger("tracing_demo", "192.168.3.36:6831")

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// ......
}
```

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		tracer.UnaryClientTracing(),
	))

	return options
}

func main() {
	// 连接jaeger服务端
	_, err := tracer.InitJaeger("hello_server", "192.168.3.36:6831")
	if err != nil {
		panic(err)
	}

	conn, _ := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

	// ......
}
```

<br>

#### gin tracing

```go
func main() {
	// 连接jaeger
	closer, _ := tracer.InitJaeger(serviceName, agentAddr)
	defer closer.Close()

	// 连接rpc服务端
	connectRPCServer(rpcAddr)

	r := gin.Default()
	r.Use(tracer.GinMiddleware()) // 添加链路跟踪中间件
	r.POST("/hello", sayHello)

	r.Run(webAddr)
}
```