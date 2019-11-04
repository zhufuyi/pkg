package npub

import (
	"strings"

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

// PushJSON 推送json数据，注：参数v必须为引用或指针
func (n *Client) PushJSON(topic string, v interface{}) error {
	eConn, err := nats.NewEncodedConn(n.Conn, nats.JSON_ENCODER)
	if err != nil {
		return err
	}

	return eConn.Publish(topic, v)
}

// PushString 推送字符串数据
func (n *Client) PushString(topic string, msg []byte) error {
	return n.Conn.Publish(topic, msg)
}

// GetClient 获取nats操作对象
func GetClient() *Client {
	if client == nil {
		panic("nats conn is nil, please connect nats server.")
	}

	return client
}
