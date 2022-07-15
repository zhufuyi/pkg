## snowflake

在[goSnowFlake](https://github.com/zheng-ji/goSnowFlake)基础上封装。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/snowflake

<br>

## 使用

```go
// 初始化
Init(1)

// 生产id
id := NewID()
```

压测：

```bash
os: windows
goarch: amd64
cpu: Intel(R) Core(TM) i7-8700 CPU @ 3.20GHz
BenchmarkNewID-12       47746561                24.72 ns/op            0 B/op          0 allocs/op
PASS
ok      command-line-arguments  1.261s
```