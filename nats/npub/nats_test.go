package npub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/zhufuyi/pkg/utils"
)

var natsAddr = []string{"nats://192.168.101.88:4222"}

func init() {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		if err := Init(natsAddr); err != nil {
			panic(err)
		}
	})
}

func TestClient_PushString(t *testing.T) {
	defer func() { recover() }()

	msg := fmt.Sprintf("hello-%d", time.Now().Unix())
	err := GetClient().PushString("foo.string", []byte(msg))
	if err != nil {
		t.Error(err)
	}
}

func TestClient_PushJSON(t *testing.T) {
	defer func() { recover() }()

	msg := struct {
		Name     string `json:"name"`
		Gender   string `json:"gender"`
		Birthday string `json:"birthday"`
	}{"张三", "男", time.Now().AddDate(-10, 0, 0).Format("2006-01-02T15:04:05Z")}

	err := GetClient().PushJSON("foo.json", &msg)
	if err != nil {
		t.Error(err)
	}
}
