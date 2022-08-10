## middleware

### 使用示例

#### jwt

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// token鉴权
	options = append(options, grpc.UnaryInterceptor(middleware.UnaryServerJwtAuth()))

	return options
}

func main() {
	fmt.Println("start rpc server", grpcAddr)
	middleware.AddSkipMethods("/proto.Account/Register") // 添加忽略token验证的方法，从pb文件的fullMethodName

	listen, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)

    // ......
}
```

<br>

#### logging

```go
var logger *zap.Logger

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	//middleware.AddLoggingFields(map[string]interface{}{"hello": "world"}) // 添加打印自定义字段
	//middleware.AddSkipLoggingMethods("/proto.Greeter/SayHello") // 跳过打印调用的方法
	options = append(options, grpc_middleware.WithUnaryServerChain(
		middleware.UnaryServerCtxTags(),
		middleware.UnaryServerZapLogging(logger),
	))

	return options
}

func main() {
	logger, _ = zap.NewProduction()

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// ......
}
```

<br>

#### recovery

```go
func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	recoveryOption := grpc_middleware.WithUnaryServerChain(
		middleware.UnaryServerRecovery(),
	)
	options = append(options, recoveryOption)

	return options
}

func main() {
	logger, _ = zap.NewProduction()

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// ......
}
```

<br>

#### retry

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 重试
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryClientRetry(
                //middleware.WithRetryTimes(5), // 修改默认重试次数，默认3次
                //middleware.WithRetryInterval(100*time.Millisecond), // 修改默认重试时间间隔，默认50毫秒
                //middleware.WithRetryErrCodes(), // 添加触发重试错误码，默认codes.Internal, codes.DeadlineExceeded, codes.Unavailable
			),
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	client := pb.NewGreeterClient(conn)

	// ......
}
```

<br>

#### timeout

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 超时拦截器
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.ContextTimeout(),
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, _ := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

    // ......
}
```
