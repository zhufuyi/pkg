# pkg

开发过程中常使用pkg库。

## 安装

> go get -u github.com/zhufuyi/pkg

<br>

## pkg列表

- [app 优雅的启动和停止服务](app)
- [awss3 aws s3客户端](awss3)
- [cache 内存和redis缓存](cache)
- [dingtalk 钉钉机器人客户端](dingtalk)
- [email 发邮件客户端](email)
- [encoding json或proto编解码](encoding)
- [errcode http和grpc错误码](errcode)
- [gin 相关](gin)
  - [validator gin请求参数校验](gin/validator)
  - [response gin返回数据封装](gin/response)
  - [errorcode 错误码定义](gin/errcode)
  - [metrics 监控指标](gin/middleware/metrics)
  - [ratelimiter 限流](gin/middleware/ratelimiter)
  - [middleware gin中间件](gin/middleware)
- [grpc 相关](grpc)
  - [benchmark 压测](grpc/benchmark)
  - [etcd grpc服务注册与发现](grpc/etcd/discovery)
  - [gtls TLS加密传输](grpc/gtls)
  - [hystrix 熔断](grpc/hystrix)
  - [keepalive 保持连接](grpc/keepalive)
  - [loadbalance 负载均衡](grpc/loadbalance)
  - [metrics grpc指标](grpc/metrics)
  - [tracer 链路跟踪](grpc/tracer)
  - [middleware 一些grpc中间件jwt、logging、recovery、retry、tracing、timeout](grpc/middleware)
- [gobash bash命令](gobash)
- [gocron 定时任务](gocron)
- [gocrypto 加密解密](gocrypto)
- [gofile 文件处理](gofile)
- [gohttp http客户端](gohttp)
- [goredis redis客户端](goredis)
- [jwt 鉴权](jwt)
- [jy2struct json或yaml转struct](jy2struct)
- [krand 随机数和字符串生成器](krand)
- [logger 日志](logger)
- [mconf 文本处理](mconf)
- [mongo 客户端](mongo)
- [mysql 客户端](mysql)
- [nats 客户端](nats)
- [redis 客户端](redis)
- [replacer 替换模板内容](replacer)
- [snowflake id生成器](snowflake)
- [sql2code 根据sql生成不同用途代码](sql2code)
- [tracer 链路跟踪](tracer)
- [utils](utils)
