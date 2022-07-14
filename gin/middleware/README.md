## render

gin中间件插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/middleware

<br>

## 使用

### 日志中间件

可以设置打印最大长度、添加请求id字段、忽略打印path、自定义zap log

```go
r := gin.Default()

// 默认打印日志
r.Use(Logging())

// 自定义打印日志
r.Use(Logging(
    WithMaxLen(400),
    WithRequestID(),
    //WithIgnoreRoutes("/hello"), // 忽略/hello
))

// 自定义zap log
log, _ := logger.Init(logger.WithFormat("json"))
r.Use(Logging(
    WithLog(log),
))
```

<br>

### 允许跨域请求

```go
r.Use(Cors())
```

<br>

### qps限制

#### 路径维度的qps限制

```go
// path, qps=500, burst=1000
r.Use(QPS())

// path, qps=20, burst=100
r.Use(ratelimiter.QPS(
    ratelimiter.WithQPS(20),
    ratelimiter.WithBurst(100),
))
```

#### ip维度的qps限制

```go
// ip, qps=10, burst=100
r.Use(QPS(
    WithIP(),
    WithQPS(10),
    WithBurst(100),
))
```