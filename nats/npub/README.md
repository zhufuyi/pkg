## npub

nats的推送消息客户端。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/nats/npub

<br>

## 使用示例

```go
    var natsAddr = []string{"nats://192.168.101.88:4222"}

	err := npub.Init(natsAddr)  // 连接

	err := npub.GetClient().PushString("foo.string", []byte(msg))  // 发送字符串

	err := npub.GetClient().PushJSON("foo.json", &msg)   // 发送对象
```
