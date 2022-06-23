## render

gin中间件插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/middleware

<br>

## 使用

日志中间件

```go
r.Use(InOutLog())
// 或
r.Use(InOutLog(
    WithMaxLen(500),
    WithIgnoreRoutes("/ping","/pong")
	))
```

<br>

允许跨域请求

```go
r.Use(Cors())
```