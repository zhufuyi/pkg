package dingTalk

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func init() {
	var (
		name        = "robot1"
		accessToken = "9146327e87b5232ceccaab36d7858f327960a335e0798920eaca8a6cb4c4f5f1"
		secret      = "SEC153e828c1294afd15771cbe23136f92e81351b5773c0842b6dd002ed520b561e"
	)
	Init([]TokenSecret{{name, accessToken, secret}})
}

func TestSendTextMessage(t *testing.T) {
	client, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	msg := NewTextMessage().SetContent("测试文本 & @某个人").SetAt([]string{"168xxxxxx"}, false)
	resp, err := client.Send(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*resp)
}

func TestSendLinkMessage(t *testing.T) {
	client, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	msg := NewLinkMessage().SetLink(
		"链接消息测试title",
		"测试text",
		"https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18b",
		"https://www.baidu.com/")
	resp, err := client.Send(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*resp)
}

func TestSendMarkdownMessage(t *testing.T) {
	client, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	mdText := `
# 标题1
## 标题2
### 标题3
#### 标题4

![screenshot](https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18bg)

###### 11点00分发布 [天气](http://www.thinkpage.cn/)
`
	msg := NewMarkdownMessage().SetMarkdown("markdown消息测试title", mdText).SetAt([]string{"135xxxxxx"}, false)
	resp, err := client.Send(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*resp)
}

func TestSendFeedCardMessage(t *testing.T) {
	client, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	msg := NewFeedCardMessage().AppendLink(
		"链接feedCard消息测试title",
		"https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18b",
		"https://www.baidu.com/")
	resp, err := client.Send(msg)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*resp)
}

func TestCheckFrequency(t *testing.T) {
	tss := new([limitFrequency]int)
	for i := 0; i < 120; i++ {
		fmt.Println(checkFrequency(tss, int(time.Now().Unix())))
		time.Sleep(time.Second)
	}
}

func TestInitDingTalk(t *testing.T) {
	Init([]TokenSecret{
		{"robot1", "token1", "secret1"},
		{"robot2", "token2", "secret2"},
		{"robot3", "token3", "secret3"},
	})

	for i := 0; i < 300; i++ {
		time.Sleep(time.Second)
		client, err := Get()
		if err != nil {
			t.Error(err)
			continue
		}
		fmt.Println(client.AccessToken)
	}
}

func TestConcurrentGet(t *testing.T) {
	Init([]TokenSecret{
		{"robot1", "token1", "secret1"},
		{"robot2", "token2", "secret2"},
		{"robot3", "token3", "secret3"},
		{"robot4", "token4", "secret4"},
		{"robot5", "token5", "secret5"},
		{"robot6", "token6", "secret6"},
	})

	wg := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, err := Get()
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}

	wg.Wait()
}

func BenchmarkGet(b *testing.B) {
	Init([]TokenSecret{
		{"robot1", "token1", "secret1"},
		{"robot2", "token2", "secret2"},
		{"robot3", "token3", "secret3"},
		{"robot4", "token4", "secret4"},
		{"robot5", "token5", "secret5"},
		{"robot6", "token6", "secret6"},
	})
	for i := 0; i < b.N; i++ {
		Get()
	}
}
