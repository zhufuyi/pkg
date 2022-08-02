# dingTalk

钉钉机器人发消息库，统一管理多个机器人，支持并发安全和限速控制，每个机器人发送限制频率20/分，如果使用6个机器人发送频率为120/分。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/dingtalk

<br>

## 使用示例

```go
    // 初始化
    dingtalk.Init([]dingtalk.TokenSecret{
        {name1, accessToken1, secret1},
        {name2, accessToken2, secret2},
        {name3, accessToken3, secret3},
    })


    // 发送文本信息
	msg := dingtalk.NewTextMessage().SetContent("测试文本 & @某个人").SetAt([]string{"168xxxxxx"}, false)
	client, err := dingtalk.Get()
	_, resp, err := client.Send(msg)


	// 发送link信息
	msg := dingtalk.NewLinkMessage().SetLink(
		"链接消息测试title",
		"测试text",
		"https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18b",
		"https://www.baidu.com/")
	client, err := dingtalk.Get()
	_, resp, err := client.Send(msg)


	// 发送markdown信息
	mdText := `
# 标题1
## 标题2
### 标题3
#### 标题4
![screenshot](https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18bg)
###### 11点00分发布 [天气](http://www.thinkpage.cn/)
`
	msg := dingtalk.NewMarkdownMessage().SetMarkdown("markdown消息测试title", mdText).SetAt([]string{"135xxxxxx"}, false)
    client, err := dingtalk.Get()
	_, resp, err := client.Send(msg)


    // 发送feedCard消息
	msg := dingtalk.NewFeedCardMessage().AppendLink(
		"链接feedCard消息测试title",
		"https://pic3.zhimg.com/v2-8962d626fed273e01f1ad08ebddf4ed5_1440w.jpg?source=172ae18b",
		"https://www.baidu.com/")
	client, err := dingtalk.Get()
	_, resp, err := client.Send(msg)

```

api接口参考文档：https://www.cnblogs.com/tjp40922/p/11299023.html
