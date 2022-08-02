package gobash

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func init() {
	if runtime.GOOS == "windows" {
		SetExecutorPath("D:\\Program Files\\cmder\\vendor\\git-for-windows\\bin\\bash.exe")
	}
}

func TestExecCombined(t *testing.T) {
	cmds := []string{
		"pwd",
		"myErrorCommand",
		"for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done",
	}

	for _, cmd := range cmds {
		out, err := ExecCombined(cmd)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Log(string(out))
	}
}

func TestExecCommand(t *testing.T) {
	cmds := []string{
		"pwd",
		"myErrorCommand",
		"for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done",
	}

	for _, cmd := range cmds {
		out, err := ExecCommand(cmd)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Log(string(out))
	}
}

func TestExecBlock(t *testing.T) {
	command := "for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done"

	result := &Result{}
	ExecRealtime(command, result)
	for v := range result.StdOut { // 通道关闭会结束for循环
		fmt.Printf(v)
	}
	if result.Err != nil {
		fmt.Println("exec command failed，", result.Err.Error())
	}
}

func TestExec(t *testing.T) {
	command := "for i in $(seq 1 5); do echo 'test cmd' $i;sleep 1; done"

	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) // 超时控制
	result := &Result{}
	// 执行
	Exec(ctx, command, result)
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
}
