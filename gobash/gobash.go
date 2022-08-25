package gobash

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// linux default path
var executor = "/bin/bash"

// SetExecutorPath 设置执行器的路径
func SetExecutorPath(path string) {
	executor = path
}

// ExecCombined 适合执行单条非阻塞命令，但是没有输出错误日志，只输出错误码，日志输出不是实时，注：当执行命令永久阻塞，会造成协程泄露
func ExecCombined(command string) ([]byte, error) {
	// 生成cmd命令
	cmd := exec.Command(executor, "-c", command)

	// 执行cmd，捕获子进程的标准输出日志
	return cmd.CombinedOutput()
}

// ExecCommand 适合执行单条非阻塞命令，输出标准和错误日志，但日志输出不是实时，注：当执行命令永久阻塞，会造成协程泄露
func ExecCommand(command string) ([]byte, error) {
	cmd := exec.Command(executor, "-c", command)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		if len(bytesErr) != 0 {
			return nil, errors.New(string(bytesErr))
		}
		return nil, err
	}

	return bytes, nil
}

// ExecRealtime 执行命令，不能主动结束命令，执行结果实时返回在channel中，注：执行命令永久阻塞，会造成协程泄露
func ExecRealtime(command string, result *Result) {
	initResult(result)

	go func() {
		defer func() { close(result.StdOut) }() // 执行完毕，关闭通道

		cmd := exec.Command(executor, "-c", command)
		handleExec(cmd, result)
	}()
}

// Exec 执行命令，可以主动结束命令，执行结果实时返回在channel中
func Exec(ctx context.Context, command string, result *Result) {
	initResult(result)

	go func() {
		defer func() { close(result.StdOut) }() // 执行完毕，关闭通道

		cmd := exec.CommandContext(ctx, executor, "-c", command)
		handleExec(cmd, result)
	}()
}

// Result 执行命令的结果
type Result struct {
	StdOut chan string
	Err    error
}

func initResult(result *Result) {
	if result == nil {
		result = &Result{StdOut: make(chan string), Err: error(nil)}
		return
	}

	if result.StdOut == nil {
		result.StdOut = make(chan string)
	}
}

func handleExec(cmd *exec.Cmd, result *Result) {
	result.StdOut <- strings.Join(cmd.Args, " ") + "\n"

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		result.Err = fmt.Errorf("stdout error, err = %s", err.Error())
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		result.Err = fmt.Errorf("stderr error, err = %s", err.Error())
		return
	}

	err = cmd.Start()
	if err != nil {
		result.Err = fmt.Errorf("cmd start error, err = %s", err.Error())
		return
	}

	reader := bufio.NewReader(stdout)
	// 实时读取每行内容
	line := ""
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			// 判断是否已经读取完毕
			if err.Error() == io.EOF.Error() {
				break
			}

			result.Err = fmt.Errorf("stdout error, err = %s", err.Error())
			break
		}
		result.StdOut <- line
	}

	// 捕获错误日志
	bytesErr, err := io.ReadAll(stderr)
	if err != nil {
		result.Err = fmt.Errorf("read stderr error, err = %s", err.Error())
		return
	}

	err = cmd.Wait()
	if err != nil {
		if len(bytesErr) != 0 {
			result.Err = fmt.Errorf("%s", bytesErr)
			return
		}
		result.Err = fmt.Errorf("cmd wait error, err = %s", err.Error())
	}
}
