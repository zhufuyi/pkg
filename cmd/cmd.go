package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	// Stdout 输出到终端
	Stdout = 1
	// File 输出到文件
	File = 2
)

var (
	filePath = "out.log" // 输出文件名和路径
	codeType = "utf8"    // codeType 输出编码类型
)

// SetFilePath 设置自定输出文件名
func SetFilePath(file string) error {
	if file == "" {
		return errors.New("file name is empty")
	}
	filePath = file
	return nil
}

// SetGbkToUtf8 设置编码类型为gbk转utf8
func SetGbkToUtf8() {
	codeType = "gbk"
}

// GbkToUtf8 把编码gbk转为utf8
func GbkToUtf8(s []byte) []byte {
	bs := []byte{}
	err := error(nil)

	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	bs, err = ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("ioutil.ReadAll error", err)
	}

	return bs
}

// 通过管道同步获取日志的函数
func syncLog(logger *log.Logger, reader io.ReadCloser) {
	//因为logger的print方法会自动添加一个换行，所以我们需要一个cache暂存不满一行的log
	cache := ""
	buf := make([]byte, 2048, 2048)
	for {
		strNum, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "file already closed") {
				return
			}
		}
		if strNum > 0 {
			outputByte := buf[:strNum]
			// 如果命令执行返回的是gbk编码，需要转换为utf8后，才不会出现乱码
			if codeType == "gbk" {
				outputByte = GbkToUtf8(outputByte)
			}
			// 这里的切分是为了将整行的log提取出来，然后将不满整行和下次一同打印
			outputSlice := strings.Split(string(outputByte), "\n")
			logText := strings.Join(outputSlice[:len(outputSlice)-1], "\n")
			logger.Printf("%s%s", cache, logText)
			cache = outputSlice[len(outputSlice)-1]
		}
	}
}

// Exec 执行命令
func Exec(command string, outType ...int) error {
	var (
		writer io.Writer
		f      *os.File
		err    error
	)
	// 设置输出类型，默认输出到终端，只有指定outType值，才根据第一个参数值来解析具体输出类型
	if len(outType) > 0 { // 只有第一个参数有效
		if outType[0] == File || outType[0] == (Stdout|File) {
			f, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
			defer f.Close()
		}

		switch outType[0] {
		case Stdout: // 只输出到终端
			writer = os.Stdout
		case File: // 只输出到文件
			writer = f
		case Stdout | File: // 同时输出到文件和终端
			writer = io.MultiWriter([]io.Writer{f, os.Stdout}...)
		default: // 默认只输出到终端
			writer = os.Stdout
		}
	} else {
		writer = os.Stdout
	}

	logger := log.New(writer, "", log.LstdFlags)
	oldFlags := logger.Flags()
	if len(command) > 100 {
		logger.Printf("run command: (%s)\n", command[:100] + " ...... ")
	} else {
		logger.Printf("run command: (%s)\n", command)
	}
	logger.SetFlags(0) // 关闭logger自身的格式，保证shell的输出和标准的log格式不冲突

	// 兼容不同OS
	name := "bash"
	arg := "-c"
	if runtime.GOOS == "windows" {
		name = "cmd"
		arg = "/C"
	}

	cmd := exec.Command(name, arg, command)
	// 获取标准输出和标准错误输出的两个管道
	cmdStdoutPipe, _ := cmd.StdoutPipe()
	cmdStderrPipe, _ := cmd.StderrPipe()
	go syncLog(logger, cmdStdoutPipe)
	go syncLog(logger, cmdStderrPipe)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()
	if err != nil {
		return err
	}

	logger.SetFlags(oldFlags) // 执行完后再打开log输出的格式
	return nil
}
