## errcode

错误码通常包括系统级错误码和服务级错误码，一共5位十进制数字组成，例如20101

| 2 | 01 | 01 |
| :------ | :------ | :------ |
| 2表示服务级错误(1为系统级错误) | 服务模块代码 | 具体错误代码 |

- 错误级别占一位数：1表示系统级错误，2表示服务级错误，通常是由用户非法操作引起的。
- 服务模块占两位数：一个大型系统的服务模块通常不超过两位数，如果超过，说明这个系统该拆分了。
- 错误码占两位数：防止一个模块定制过多的错误码，后期不好维护。

<br>

### 安装

> go get -u github.com/zhufuyi/pkg/errcode

<br>

### 使用示例

### http错误码使用示例

```go
    // 定义错误码
    var ErrLogin = errcode.NewError(20101, "用户名或密码错误")

    // 请求返回
    response.Error(c, errcode.LoginErr)
```

<br>

### grpc错误码使用示例

```go
    // 定义错误码
    var ErrLogin = errcode.NewRPCErr(20101, "用户名或密码错误")

    // 返回错误
    // req *pb.CreateRequest
    errcode.RPCErr(req, errcode.ErrLogin)
    // errcode.RPCErr(req, errcode.ErrLogin, errcode.KV{"msg":err.Error()}) // 附带错误详情信息
```