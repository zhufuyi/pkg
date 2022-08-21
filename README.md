# pkg

开发过程中常使用pkg库。

## 安装

> go get -u github.com/zhufuyi/pkg

<br>

## pkg列表

- [awss3 aws s3客户端](awss3)
- [cache 内存和redis缓存](cache)
- [dingtalk 钉钉机器人客户端](dingtalk)
- [email 发邮件客户端](email)
- [encoding json或proto编解码](encoding)
- [gin 相关](gin)
  - [validator gin请求参数校验](gin/validator)
  - [response gin返回数据封装](gin/render)
  - [errorcode 错误码定义](gin/errcode)
  - [middleware gin中间件](gin/middleware)
- [grpc 相关](grpc)
  - [errcode grpc错误码](grpc/errcode)
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
- [krand 随机数和字符串生成器](krand)
- [logger 日志](logger)
- [mconf 文本处理](mconf)
- [mongo 客户端](mongo)
- [mysql 客户端](mysql)
- [nats 客户端](nats)
- [redis 客户端](redis)
- [snowflake id生成器](snowflake)
- [tracer 链路跟踪](tracer)
