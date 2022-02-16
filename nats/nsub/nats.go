package nsub

import (
	"strings"
	"context"
	"time"
	"fmt"

	"github.com/nats-io/nats.go"
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
func (c *Client) Close() {
	if c.Conn != nil {
		c.Conn.Close()
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
			return
		}

		select {
		case pushData <- msg.Data:
			continue
		case <-ctx.Done():
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
