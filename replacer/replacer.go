package replacer

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zhufuyi/pkg/gofile"
)

var _ Replacer = (*replacerInfo)(nil)

// Replacer 接口
type Replacer interface {
	SetReplacementFields(fields []Field)
	SetIgnoreFiles(filenames ...string)
	SetOutPath(absPath string, name ...string) error
	GetOutPath() string
	SaveFiles() error
	ReadFile(filename string) ([]byte, error)
}

// replacerInfo replace设置信息
type replacerInfo struct {
	path              string   // 模板目录路径(不包含.或..)
	fs                embed.FS // 模板目录对应二进制对象
	isActual          bool     // fs字段是否来源实际路径，如果为true，使用io操作文件，如果为false使用fs操作文件
	files             []string // 模板文件列表
	ignoreFiles       []string // 忽略替换的文件列表
	replacementFields []Field  // 从模板文件转为新文件需要替换的字符
	outPath           string   // 输出替换后文件存放目录路径
}

// New 根据指定路径创建replacer
func New(path string) (*replacerInfo, error) {
	files, err := gofile.ListFiles(path)
	if err != nil {
		return nil, err
	}

	path, _ = filepath.Abs(path)
	return &replacerInfo{
		path:              path,
		isActual:          true,
		files:             files,
		replacementFields: []Field{},
	}, nil
}

// NewWithFS 根据嵌入的路径创建replacer
func NewWithFS(path string, fs embed.FS) (Replacer, error) {
	files, err := listFiles(path, fs)
	if err != nil {
		return nil, err
	}

	return &replacerInfo{
		path:              path,
		fs:                fs,
		isActual:          false,
		files:             files,
		replacementFields: []Field{},
	}, nil
}

// Field 替换字段信息
type Field struct {
	Old             string // 模板字段
	New             string // 新字段
	IsCaseSensitive bool   // 第一个字母是否区分大小写
}

// SetReplacementFields 设置替换字段，注：old字符尽量不要存在包含关系，如果存在，在设置Field时注意先后顺序
func (t *replacerInfo) SetReplacementFields(fields []Field) {
	var newFields []Field
	for _, field := range fields {
		if field.IsCaseSensitive && isFirstAlphabet(field.Old) { // 拆分首字母大小写两个字段
			newFields = append(newFields,
				Field{ // 把第一个字母转为大写
					Old: strings.ToUpper(field.Old[:1]) + field.Old[1:],
					New: strings.ToUpper(field.New[:1]) + field.New[1:],
				},
				Field{ // 把第一个字母转为小写
					Old: strings.ToLower(field.Old[:1]) + field.Old[1:],
					New: strings.ToLower(field.New[:1]) + field.New[1:],
				},
			)
		} else {
			newFields = append(newFields, field)
		}
	}
	t.replacementFields = newFields
}

// SetIgnoreFiles 设置忽略处理的文件
func (t *replacerInfo) SetIgnoreFiles(filenames ...string) {
	t.ignoreFiles = append(t.ignoreFiles, filenames...)
}

// SetOutPath 设置输出目录路径，优先使用absPath绝对路径，如果absPath为空，自动在当前目录根据参数name和时间生成绝对路径
func (t *replacerInfo) SetOutPath(absPath string, name ...string) error {
	subPath := ""
	if len(name) > 0 && name[0] != "" {
		subPath = name[0]
	}

	if absPath != "" {
		abs, err := filepath.Abs(absPath)
		if err != nil {
			return err
		}

		t.outPath = abs + gofile.GetPathDelimiter() + subPath
		return nil
	}

	t.outPath = gofile.GetRunPath() + gofile.GetPathDelimiter() + subPath + "_" + time.Now().Format("0102150405")
	return nil
}

// GetOutPath 获取输出目录路径
func (t *replacerInfo) GetOutPath() string {
	return t.outPath
}

// ReadFile 读取文件内容
func (t *replacerInfo) ReadFile(filename string) ([]byte, error) {
	count := 0
	foundFile := ""
	for _, file := range t.files {
		if strings.Contains(file, filename) {
			count++
			foundFile = file
		}
	}
	if count == 0 || count > 1 {
		return nil, fmt.Errorf("total %d file named '%s'", count, filename)
	}

	if t.isActual {
		return os.ReadFile(foundFile)
	}
	return t.fs.ReadFile(foundFile)
}

// SaveFiles 导出文件
func (t *replacerInfo) SaveFiles() error {
	if t.outPath == "" {
		t.outPath = gofile.GetRunPath() + gofile.GetPathDelimiter() + "template_" + time.Now().Format("0102150405")
	}

	for _, file := range t.files {
		if t.isIgnoreFile(file) {
			continue
		}

		// 从二进制读取模板文件内容使用embed.FS，如果要从指定目录读取使用os.ReadFile
		var data []byte
		var err error
		if t.isActual {
			data, err = os.ReadFile(file)
		} else {
			data, err = t.fs.ReadFile(file)
		}
		if err != nil {
			return err
		}

		// 替换文本内容
		for _, field := range t.replacementFields {
			data = bytes.ReplaceAll(data, []byte(field.Old), []byte(field.New))
		}

		// 获取新文件路径
		newFilePath := t.getNewFilePath(file)
		dir, filename := filepath.Split(newFilePath)
		// 替换文件名
		for _, field := range t.replacementFields {
			tmp := dir + strings.ReplaceAll(filename, field.Old, field.New)
			if newFilePath != tmp {
				newFilePath = tmp
				break
			}
		}

		// 保存文件
		err = saveToNewFile(newFilePath, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *replacerInfo) isIgnoreFile(file string) bool {
	isIgnore := false
	_, filename := filepath.Split(file)
	for _, v := range t.ignoreFiles {
		if filename == v {
			isIgnore = true
			break
		}
	}
	return isIgnore
}

func (t *replacerInfo) getNewFilePath(file string) string {
	var newFilePath string
	if t.isActual {
		newFilePath = t.outPath + strings.Replace(file, t.path, "", 1)
	} else {
		newFilePath = t.outPath + strings.Replace(file, t.path, "", 1)
	}

	if runtime.GOOS == "windows" {
		newFilePath = strings.ReplaceAll(newFilePath, "/", "\\")
	}

	return newFilePath
}

func saveToNewFile(filePath string, data []byte) error {
	// 创建目录
	dir, _ := filepath.Split(filePath)
	err := os.MkdirAll(dir, 0666)
	if err != nil {
		return err
	}

	// 保存文件
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// 遍历嵌入的目录下所有文件，返回文件的绝对路径
func listFiles(path string, fs embed.FS) ([]string, error) {
	files := []string{}
	err := walkDir(path, &files, fs)
	return files, err
}

// 通过迭代方式遍历嵌入的目录
func walkDir(dirPath string, allFiles *[]string, fs embed.FS) error {
	files, err := fs.ReadDir(dirPath) // 读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + "/" + file.Name()
		if file.IsDir() {
			walkDir(deepFile, allFiles, fs)
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

// 判断字符串第一个字符是字母
func isFirstAlphabet(str string) bool {
	if len(str) == 0 {
		return false
	}

	if (str[0] >= 'A' && str[0] <= 'Z') || (str[0] >= 'a' && str[0] <= 'z') {
		return true
	}

	return false
}
