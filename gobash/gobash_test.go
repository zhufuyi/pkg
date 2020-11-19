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
		SetExecutorPath("C:\\tools\\cmder\\vendor\\git-for-windows\\usr\\bin\\bash.exe")
	}
}

func TestSimpleCMD(t *testing.T) {
	cmds := []string{
		"pwd",
		"myErrorCommand",
		"for i in $(seq 1 3); do echo 'test cmd' $i;sleep 0.5; done",
	}

	for _, cmd := range cmds {
		out, err := Simple(cmd)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(string(out))
		}
	}
}

func TestExecNBCMD(t *testing.T) {
	cmds := []string{
		"pwd",
		"myErrorCommand",
		"for i in $(seq 1 3); do echo 'test cmd' $i;sleep 0.2; done",
	}

	for _, cmd := range cmds {
		std, err := ExecNonBlock(cmd)
		if err != nil {
			fmt.Println(err)
		} else {
			if std != nil {
				fmt.Println(string(std))
			}
		}
	}
}

func TestExecBlockCMD(t *testing.T) {
	command := "for i in $(seq 1 3); do echo 'test cmd' $i;sleep 1; done"

	result := &Result{}
	ExecBlock(command, result)
	for v := range result.StdOut { // 通道关闭会结束for循环
		fmt.Printf(v)
	}
	if result.Err != nil {
		fmt.Println("exec command failed，", result.Err.Error())
	}
}

func TestExecCMD(t *testing.T) {
	// 测试主动结束命令
	command := "for i in $(seq 1 10); do echo 'test cmd' $i;sleep 1; done"

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := &Result{}

	Exec(ctx, command, result)
	for v := range result.StdOut { // 通道关闭会自动结束for循环
		fmt.Printf(v)
	}
	if result.Err != nil {
		if ctx.Err() != nil {
			fmt.Println("kill the cmd success")
			return
		}
		fmt.Println(ctx.Err(), "exec command failed,", result.Err.Error())
	}
}
