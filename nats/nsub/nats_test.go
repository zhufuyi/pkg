package nsub

import (
	"context"
	"fmt"
	"testing"
	"time"

	"pkg/nats/npub"
)

var natsAddr = []string{"nats://192.168.101.88:4222"}

func init() {
	if err := npub.Init(natsAddr); err != nil {
		panic(err)
	}

	if err := Init(natsAddr); err != nil {
		panic(err)
	}
}

func TestClient_SubscribeSync(t *testing.T) {
	topic := "foo.json"
	pushData := make(chan []byte, 100)
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)

	// 推送
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				msg := struct {
					Name     string `json:"name"`
					Gender   string `json:"gender"`
					Birthday string `json:"birthday"`
				}{"张三", "男", time.Now().AddDate(-10, 0, 0).Format("2006-01-02T15:04:05Z")}

				err := npub.GetClient().PushJSON(topic, &msg)
				if err != nil {
					t.Error(err)
				}
				fmt.Printf("[pub] %s\n", msg)
			case <-ctx.Done():
				return
			}
		}
	}()

	// 订阅
	go func() {
		GetClient().SubscribeSync(ctx, topic, pushData)
	}()

	for {
		select {
		case msg := <-pushData:
			fmt.Printf("[sub] %s\n\n", msg)
		case <-ctx.Done():
			return
		}
	}
}
