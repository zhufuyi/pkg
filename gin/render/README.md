## render

gin请求和返回的封装插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/render

<br>

## 使用示例

`Respond`函数返回兼容http状态码

`Success`和`Error`统一返回状态码200，在data.code自定义状态码

所有请求统一返回json

```json
{
  "code": 0,
  "msg": "",
  "data": {}
}
```

```go
    // c是*gin.Context

    // 返回成功
    render.Success(c)
    // 返回成功，并返回数据
    render.Success(c, gin.H{"users":users})

    // 返回失败
    render.Error(c, errcode.SendEmailErr)
    // 返回失败，并返回数据
    render.Error(c,  errcode.SendEmailErr, gin.H{"user":user})
```