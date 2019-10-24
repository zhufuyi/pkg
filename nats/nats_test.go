package nats

import (
	"fmt"
	"testing"
	"time"
)

var natsAddr = []string{"nats://192.168.101.88:4222"}

func init() {
	if err := Init(natsAddr); err != nil {
		panic(err)
	}
}

func TestClient_PushString(t *testing.T) {
	msg := fmt.Sprintf("hello-%d", time.Now().Unix())
	err := GetClient().PushString("foo.string", []byte(msg))
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 100)
}

func TestClient_PushJSON(t *testing.T) {
	msg := struct {
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
	}{"张三", "男", time.Now().AddDate(-10, 0, 0).Format("2006-01-02T15:04:05Z")}

	err := GetClient().PushJSON("foo.json", &msg)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Millisecond * 100)
}
