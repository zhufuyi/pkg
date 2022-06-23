## render

gin请求和返回的封装插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/render

<br>

## 使用

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

具体使用看 [测试文件](output_test.go)
