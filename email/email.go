package email

import (
	"crypto/tls"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

// https://pkg.go.dev/gopkg.in/gomail.v2#example-package

// Client 发送客户端
type Client struct {
	Host     string // host smtp地址
	Port     int    // port 端口
	Username string // username 账户
	Password string // password 密码
}

// Init 实例化
func Init(username string, password string) (*Client, error) {
	client := &Client{
		Username: username,
		Password: password,
		Port:     465,
	}

	split := strings.Split(username, "@")
	if len(split) != 2 {
		return nil, fmt.Errorf("'%s' is not email address", username)
	}
	client.Host = "smtp." + split[1]

	return client, nil
}

// Message 内容
type Message struct {
	To          []string // 收件人
	Cc          []string // 抄送人
	Subject     string   // 标题
	ContentType string   // 内容的类型text/plain,text/html
	Content     string   // 发送内容
	Attach      string   // 附件
}

// SendMessage 发送邮件
func (c *Client) SendMessage(msg *Message) error {
	gm := gomail.NewMessage()
	gm.SetHeader("From", c.Username)
	gm.SetHeader("To", msg.To...)
	if len(msg.Cc) > 0 {
		gm.SetHeader("Cc", msg.Cc...)
	}
	gm.SetHeader("Subject", msg.Subject)
	gm.SetBody(msg.ContentType, msg.Content)
	if msg.Attach != "" {
		gm.Attach(msg.Attach)
	}

	dialer := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialer.DialAndSend(gm)
}
