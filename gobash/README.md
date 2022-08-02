## s3 api

在go环境中执行命令、脚本、可执行文件。

- 支持linux和windows系统，支持gbk转utf8
- 执行命令过程日志实时输出到终端或文件，也可以同时写入到终端和文件

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gobash

<br>

## 使用示例

(1) Exec 执行命令，可以主动结束命令，执行结果实时返回在channel中

```go
	command := "for i in $(seq 1 5); do echo 'test cmd' $i;sleep 1; done"

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) // 超时控制
	result := &gobash.Result{}
	// 执行
	gobash.Exec(ctx, command, result)
	// 实时输出日志和错误
Receive:
	for {
		select {
		case out, ok := <-result.StdOut:
			if !ok { // 通道关闭会结束for循环
				break Receive
			}
			fmt.Printf(out) // 实时日志

		case <-ctx.Done():
			result.Err = ctx.Err()
			break Receive
		}
	}

	if result.Err != nil {
		fmt.Println("exec command failed,", result.Err.Error())
	}
```

(2) ExecCombined 适合执行单条非阻塞命令，但是没有输出错误日志，只输出错误码，日志输出不是实时，注：当执行命令永久阻塞，会造成协程泄露

```go
    command := "for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done"
    out, err := gobash.ExecCombined(command)
    if err != nil {
        return
    }
    fmt.Println(string(out))
```

<br>

(3) ExecCommand 适合执行单条非阻塞命令，输出标准和错误日志，但日志输出不是实时，注：当执行命令永久阻塞，会造成协程泄露

```go
    command := "for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done"
    out, err := gobash.ExecCommand(command)
    if err != nil {
        return
    }
    fmt.Println(string(out))
```

<br>

(4) ExecRealtime 执行命令，不能主动结束命令，执行结果实时返回在channel中，注：执行命令永久阻塞，会造成协程泄露

```go
	command := "for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done"

	result := &gobash.Result{}
	gobash.ExecRealtime(command, result)
	for v := range result.StdOut { // 通道关闭会结束for循环
		fmt.Printf(v)
	}
	if result.Err != nil {
		fmt.Println("exec command failed，", result.Err.Error())
	}
```