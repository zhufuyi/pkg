package nsub

import (
	"strings"

	"context"
	"time"

	"fmt"

	"github.com/nats-io/go-nats"
)

var client *Client

// Client nats连接套接字
type Client struct {
	Conn *nats.Conn
}

// Init nats初始化
func Init(natsAddr []string) error {
	natsURL := strings.Join(natsAddr, ",")
	conn, err := nats.Connect(natsURL, nats.DontRandomize())
	if err != nil {
		return err
	}

	client = &Client{Conn: conn}
	return nil
}

// Close 关闭
func (n *Client) Close() {
	if n.Conn != nil {
		n.Conn.Close()
	}
}

// SubscribeSync 同步订阅
func (c *Client) SubscribeSync(ctx context.Context, topic string, pushData chan []byte) {
	sub, err := c.Conn.SubscribeSync(topic)
	if err != nil {
		fmt.Printf("[nats] SubscribeSync() error, topic=%s, err=%s\n", topic, err.Error())
		return
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		// 等待消息
		msg, err := sub.NextMsgWithContext(ctx)
		if err != nil {
			fmt.Printf("[nats] NextMsgWithContext() error, topic=%s, err=%s\n", topic, err.Error())
			return
		}

		select {
		case pushData <- msg.Data:
			continue
		case <-ctx.Done():
			fmt.Printf("[nats] exit subscribe, topic=%s, %s\n", topic, ctx.Err())
			return
		default: // forbid block
		}

		select {
		case <-ticker.C:
			fmt.Printf("[nats] push data to channel timeout, topic=%s\n", topic)
		default: // forbid block
		}
	}
}

// GetClient 获取nats操作对象
func GetClient() *Client {
	if client == nil {
		panic("nats conn is nil, please connect nats server.")
	}

	return client
}
