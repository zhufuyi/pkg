package cmd

import (
	"testing"
)

func TestExec(t *testing.T) {
	var err error
	var command string

	//command = "ping www.qq.com"
	//SetGbkToUtf8()
	//err = Exec(command)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	command = "ping www.baidu.com"
	SetFilePath("./ping.log")
	err = Exec(command, Stdout|File)
	if err != nil {
		t.Error(err)
		return
	}
}
