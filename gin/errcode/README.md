## errcode

errcode用来自定义错误码。

<br>

## 使用示例

```go
    // 定义错误码
    LoginErr = errcode.NewError(100501, "登录失败，用户名或密码错误")

    // 请求返回
    render.Error(c, errcode.LoginErr)
```