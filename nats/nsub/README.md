## nsub

nats的订阅消息客户端。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/nats/nsub

<br>

## 使用示例

```go
    var natsAddr = []string{"nats://192.168.101.88:4222"}

    // 初始化
    err := nsub.Init(natsAddr)
    
	topic := "foo.json"
	subData := make(chan []byte, 100)
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	// 订阅
	go func() {
		nsub.GetClient().SubscribeSync(ctx, topic, subData)
	}()

	for {
		select {
		case msg := <-subData:
			fmt.Printf("[sub] %s\n\n", msg)
		case <-ctx.Done():
			return
		}
	}
```

