## discovery

etcd 作为服务注册与发现。

<br>

### 使用示例

#### grpc server

```go
func startServer(addr string) {
	fmt.Println("start rpc server", addr)
	list, _ := net.Listen("tcp", addr)
	server := grpc.NewServer()
    // ......
}

func main() {
	etcdAddrs := []string{"192.168.3.36:2379"}
	serverName := "hello-demo"
	serverAddr:="127.0.0.1:8080"

	// 运行rpc服务
	go startServer(grpcAddr)

	// 注册服务到etcd
	etcdRegister := discovery.RegisterRPCAddr(serverName, serverAddr, etcdAddrs)

     // 自定义设置方式
	//etcdRegister := discovery.RegisterRPCAddr(serverName, serverAddr, etcdAddrs,
	//	discovery.WithTTLSeconds(5),  // 超时5秒
	//	discovery.WithLogger(zap.NewExample()), // 日志
	//	discovery.WithWeight(5), // 权重
	//)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	etcdRegister.Stop()
}
```

<br>

#### grpc client

```go
func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return options
}

func main() {
    serverName := "hello-demo"
    etcdAddrs := []string{"192.168.3.36:2379"}

	// 使用etcd服务发现
	r := discovery.NewResolver(etcdAddrs)
	// 自定义设置方式
	//r := discovery.NewResolver(etcdAddrs,
	//	discovery.WithDialTimeout(5),
	//	discovery.WithLogger(zap.NewExample()),
	//)
	resolver.Register(r)

	conn, _ := grpc.Dial("etcd:///"+serverName, getDialOptions()...)

    // ......
}
```