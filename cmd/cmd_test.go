package cmd

import (
	"testing"
)

func TestExec(t *testing.T) {
	var err error
	var command string

	// 只输出到终端
	//command = "ping www.qq.com"
	//SetGbkToUtf8()
	//err = Exec(command)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	// 执行命令日志同时输出到终端和文件
	command = "ping www.baidu.com"
	SetFilePath("./ping.log")
	err = Exec(command, Stdout|File)
	if err != nil {
		t.Error(err)
		return
	}
}
