package email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

// https://pkg.go.dev/gopkg.in/gomail.v2#example-package

type ServerType int

const (
	// qq邮箱
	ServerTypeQQ ServerType = 0
	// 126邮箱
	ServerType126 ServerType = 1
	// 163邮箱
	ServerType163 ServerType = 2
)

// Message 内容
type Message struct {
	To          []string // 收件人
	Cc          []string // 抄送人
	Subject     string   // 标题
	ContentType string   // 内容的类型text/plain,text/html
	Content     string   // 发送内容
	Attach      string   // 附件
}

// Client 发送客户端
type Client struct {
	Host     string // host smtp地址
	Port     int    // port 端口
	Username string // username 账户
	Password string // password 密码
}

// NewClient 实例化邮件客户端
func NewClient(username string, password string, st ServerType) *Client {
	client := &Client{
		Username: username,
		Password: password,
	}

	switch st {
	case ServerTypeQQ:
		client.Host = "smtp.qq.com"
		client.Port = 465
	case ServerType126:
		client.Host = "smtp.126.com"
		client.Port = 465
	case ServerType163:
		client.Host = "smtp.163.com"
		client.Port = 465
	}

	return client
}

// SendMessage 发送邮件
func (c *Client) SendMessage(msg *Message) (bool, error) {
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
	if err := dialer.DialAndSend(gm); err != nil {
		return false, err
	}
	return true, nil
}
