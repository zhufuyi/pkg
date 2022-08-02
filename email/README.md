# email

在[gomail](https://gopkg.in/gomail.v2)基础上封装的go语言发邮件库

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/email

<br>

## 使用示例

```go
    // 初始化，参数是为发送者邮箱和密码(不是登录密码)
    client, err := email.Init("xxxxxx@qq.com", "xxxxxx")

    // 发送邮件
	msg := &email.Message{
		To:          []string{"xxxxxx@qq.com"},
		Cc:          nil,
		Subject:     "title-demo",
		ContentType: "text/plain",
		Content:     "邮件内容demo-01",
		Attach:      "",
	}
	err := client.SendMessage(msg)
```

<br>

### qq邮箱设置

#### 开启邮箱服务

- 登录qq邮箱(手机号已认证)。
- 点击设置 --> 账户，找到STMP服务，点击开启，根据提示内容使用手机发短信到指定手机号码，开启完成后得到授权码，有了授权码就可以使用客户端收发邮件，[获取授权码教程](https://service.mail.qq.com/cgi-bin/help?subtype=1&&id=28&&no=1001256) 。

邮箱的服务器地址和端口(使用SSL)：

- 发送邮件服务器：smtp.qq.com，端口号465或587


<br>

### 网易邮箱设置

#### 开启邮箱服务

- 登录126邮箱(手机号已认证)。
- 点击设置 --> POP3/SMTP/IMAP，点击开启，根据提示内容使用手机发短信到指定手机号码，开启完成后得到授权码，有了授权码就可以使用客户端收发邮件，[获取授权码教程](https://help.mail.163.com/faqDetail.do?code=d7a5dc8471cd0c0e8b4b8f4f8e49998b374173cfe9171305fa1ce630d7f67ac21b87735d7227c217) 。

邮箱的服务器地址和端口(使用SSL)：

- 126邮箱：
  - 发送邮件服务器：smtp.126.com，端口号465或994
- 163邮箱：
  - 发送邮件服务器：smtp.163.com，端口号465或994

