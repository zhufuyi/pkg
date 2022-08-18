## render

gin中间件插件。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/middleware

<br>

## 使用示例

### 日志中间件

可以设置打印最大长度、添加请求id字段、忽略打印path、自定义zap log

```go
    r := gin.Default()

    // 默认打印日志
    r.Use(middleware.Logging())

    // 自定义打印日志
    r.Use(middleware.Logging(
        middleware.WithMaxLen(400),
        middleware.WithRequestID(),
        //middleware.WithIgnoreRoutes("/hello"), // 忽略/hello
    ))

    // 自定义zap log
    log, _ := logger.Init(logger.WithFormat("json"))
    r.Use(middlewareLogging(
        middleware.WithLog(log),
    ))
```

<br>

### 允许跨域请求

```go
    r := gin.Default()
    r.Use(middleware.Cors())
```

<br>

### qps限制

#### 路径维度的qps限制

```go
    r := gin.Default()

    // path, 默认qps=500, burst=1000
    r.Use(ratelimiter.QPS())

    // path, 自定义qps=20, burst=100
    r.Use(ratelimiter.QPS(
        ratelimiter.WithQPS(20),
        ratelimiter.WithBurst(100),
    ))
```

#### ip维度的qps限制

```go
    // ip, 自定义qps=10, burst=100
    r.Use(ratelimiter.QPS(
        ratelimiter.WithIP(),
        ratelimiter.WithQPS(10),
        ratelimiter.WithBurst(100),
    ))
```

<br>

### jwt鉴权

```go
    r := gin.Default()
    r.GET("/user/:id", middleware.JWT(), userFun) // 需要鉴权
```
<br>

### 链路跟踪

```go
    r := gin.Default()
    r.Use(middleware.Tracing("your-service-name"))
```