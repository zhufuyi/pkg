package email

import (
	"fmt"
	"testing"
)

func TestEmailClient_SendMessage(t *testing.T) {
	// 发送内容
	msg := &Message{
		To:          []string{"xxxxxx@qq.com"},
		Cc:          nil,
		Subject:     "title-demo",
		ContentType: "text/plain",
		Content:     "邮件内容demo-01",
		Attach:      "",
	}

	// qq邮箱发送
	client := NewClient("xxxxxx@qq.com", "xxxxxx", ServerTypeQQ)
	ok, err := client.SendMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	// 126邮箱发送
	client = NewClient("xxxxxx@126.com", "xxxxxx", ServerType126)
	ok, err = client.SendMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	// 163邮箱发送
	client = NewClient("xxxxxx@63.com", "xxxxxx", ServerType163)
	ok, err = client.SendMessage(msg)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(ok)
}
