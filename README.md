# pkg

<div align=center>

[![Go Report](https://goreportcard.com/badge/github.com/zhufuyi/pkg)](https://goreportcard.com/report/github.com/zhufuyi/pkg)
[![codecov](https://codecov.io/gh/zhufuyi/pkg/branch/main/graph/badge.svg)](https://codecov.io/gh/zhufuyi/pkg)
[![Go Reference](https://pkg.go.dev/badge/github.com/zhufuyi/pkg.svg)](https://pkg.go.dev/github.com/zhufuyi/pkg)
[![Go](https://github.com/zhufuyi/pkg/workflows/Go/badge.svg?branch=main)](https://github.com/zhufuyi/pkg/actions)
[![License: MIT](https://img.shields.io/github/license/zhufuyi/pkg)](https://img.shields.io/github/license/zhufuyi/pkg)

</div>

常使用的pkg库。

## 安装

> go get -u github.com/zhufuyi/pkg@latest

<br>

## pkg列表

- [app 优雅的启动和停止服务](app)
- [awss3 aws s3客户端](awss3)
- [cache 内存和redis缓存](cache)
- [conf 解析yaml、json、toml配置文件](conf)
- [consulcli 客户端](consulcli)
- [container 客户端](container)
- [dingtalk 钉钉机器人客户端](dingtalk)
- [email 发邮件客户端](email)
- [encoding json、proto和gob编解码](encoding)
- [errcode http和rpc错误码](errcode)
- [gin 相关](gin)
  - [handlerfunc 常用handler函数](gin/handlerfunc)
  - [middleware 中间件](gin/middleware)
    - [metrics 指标](gin/middleware/metrics)
    - [auth 鉴权](gin/middleware/auth.go)
    - [breaker 熔断器](gin/middleware/breaker.go)
    - [cors 跨域](gin/middleware/cors.go)
    - [logging 日志](gin/middleware/logging.go)
    - [ratelimit 限流](gin/middleware/ratelimit.go)
    - [request id 请求id](gin/middleware/requestid.go)
    - [tracing 链路跟踪](gin/middleware/tracing.go)
    - [ratelimit 限流](gin/middleware/ratelimit)
  - [prof go profile](gin/prof)
  - [response 返回数据封装](gin/response)
  - [swagger api文档](gin/swagger)
  - [validator 请求参数校验](gin/validator)
- [grpc 相关](grpc)
  - [benchmark 压测](grpc/benchmark)
  - [grpccli grpc 客户端](grpc/grpccli)
  - [gtls TLS加密传输](grpc/gtls)
  - [keepalive 保持连接](grpc/keepalive)
  - [loadbalance 负载均衡](grpc/loadbalance)
  - [metrics rpc指标](grpc/metrics)
  - [interceptor 客户端和服务端的拦截器](grpc/interceptor)
    - [breaker 熔断器](grpc/interceptor/breaker.go)
    - [jwtAuth 鉴权](grpc/interceptor/jwtAuth.go)
    - [logging 日志](grpc/interceptor/logging.go)
    - [metrics 指标](grpc/interceptor/metrics.go)
    - [ratelimit 重试](grpc/interceptor/ratelimit.go)
    - [recovery 恢复](grpc/interceptor/recovery.go)
    - [retry 重试](grpc/interceptor/retry.go)
    - [timeout 超时](grpc/interceptor/timeout.go)
    - [tracing 链路跟踪](grpc/interceptor/tracing.go)
- [gobash bash命令](gobash)
- [gocron 定时任务](gocron)
- [gocrypto 加密解密](gocrypto)
- [gofile 文件处理](gofile)
- [gohttp http客户端](gohttp)
- [goredis redis客户端](goredis)
- [gotest 测试库](gotest)
- [jwt 鉴权](jwt)
- [jy2struct json或yaml转struct](jy2struct)
- [krand 随机数和字符串生成器](krand)
- [logger 日志](logger)
- [mconf 文本处理](mconf)
- [mongo 客户端](mongo)
- [mysql 客户端](mysql)
- [nats 客户端](nats)
- [prof go profile](prof)
- [redis 客户端](redis)
- [replacer 替换模板内容](replacer)
- [servicerd 服务注册与发现](servicerd)
- [shield 系统和进程资源统计](shield)
- [snowflake id生成器](snowflake)
- [sql2code 根据sql生成不同用途代码](sql2code)
- [tracer 链路跟踪](tracer)
- [utils](utils)
