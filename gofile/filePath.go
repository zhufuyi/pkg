package gofile

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// IsExists 判断文件或文件夹是否存在
func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}

	return false
}

// GetRunPath 获取程序执行的绝对路径
func GetRunPath() string {
	dir, err := os.Executable()
	if err != nil {
		fmt.Println("os.Executable error.", err.Error())
		return ""
	}

	return filepath.Dir(dir)
}

// GetPathDelimiter 根据系统类型获取分隔符
func GetPathDelimiter() string {
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}

	return delimiter
}

// ListFiles 遍历指定目录下所有文件，返回文件的绝对路径
func ListFiles(dirPath string) ([]string, error) {
	files := []string{}
	err := error(nil)

	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return files, err
	}

	return files, walkDir(dirPath, &files)
}

// ListDirsAndFiles 遍历指定目录下所有子目录文件，返回文件的绝对路径
func ListDirsAndFiles(dirPath string) (map[string][]string, error) {
	df := make(map[string][]string, 2)

	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return df, err
	}

	dirs := []string{}
	files := []string{}
	err = walkDir2(dirPath, &dirs, &files)
	if err != nil {
		return df, err
	}

	df["dirs"] = dirs
	df["files"] = files

	return df, nil
}

type FilterFn func(string) bool

// ListFilesWithFilter 带过滤条件遍历指定目录下所有文件，返回绝对路径
func ListFilesWithFilter(dirPath string, filter FilterFn) ([]string, error) {
	files := []string{}
	err := error(nil)

	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		return files, err
	}

	return files, walkDirWithFilter(dirPath, &files, filter)
}

// MatchSuffix 后缀匹配
func MatchSuffix(suffixName string) FilterFn {
	return func(filename string) bool {
		if suffixName == "" {
			return false
		}

		size := len(filename) - len(suffixName)
		if size >= 0 && filename[size:] == suffixName { // 后缀
			return true
		}
		return false
	}
}

// MatchContain 包含字符串
func MatchContain(baseName string) FilterFn {
	return func(filename string) bool {
		if baseName == "" {
			return false
		}

		return strings.Contains(filename, baseName)
	}
}

// DeleteDir 删除指定目录下所有文件和目录
func DeleteDir(dirPath string) ([]string, error) {
	var errStr string
	var deleteFiles []string

	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return nil, err
	}
	if absPath == "/" || strings.Contains(absPath, ":\\") {
		return nil, fmt.Errorf("can't delete directory %s", dirPath)
	}
	
	df, err := ListDirsAndFiles(dirPath)
	if err != nil {
		return deleteFiles, err
	}

	files := df["files"]
	dirs := df["dirs"]

	// 删除文件
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			errStr += err.Error() + "/n"
			continue
		}
		deleteFiles = append(deleteFiles, file)
	}

	// 删除目录
	size := len(dirs)
	for i := size - 1; i >= 0; i-- {
		err = os.RemoveAll(dirs[i])
		if err != nil {
			errStr += err.Error() + "/n"
			continue
		}
		deleteFiles = append(deleteFiles, dirs[i])
	}

	// 删除指定目录
	err = os.RemoveAll(dirPath)
	if err != nil {
		errStr += err.Error() + "/n"
	}

	if errStr != "" {
		return deleteFiles, errors.New(errStr)
	}

	return deleteFiles, nil
}

func getLevel(dir string) int {
	if runtime.GOOS == "windows" {
		return strings.Count(dir, "\\")
	}
	return strings.Count(dir, "/")
}

// 通过迭代方式遍历文件
func walkDir(dirPath string, allFiles *[]string) error {
	files, err := ioutil.ReadDir(dirPath) // 读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			walkDir(deepFile, allFiles)
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

func walkDir2(dirPath string, allDirs *[]string, allFiles *[]string) error {
	files, err := ioutil.ReadDir(dirPath) // 读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			*allDirs = append(*allDirs, deepFile)
			walkDir2(deepFile, allDirs, allFiles)
			continue
		}
		*allFiles = append(*allFiles, deepFile)
	}

	return nil
}

// 带过滤条件通过迭代方式遍历文件
func walkDirWithFilter(dirPath string, allFiles *[]string, filter func(string) bool) error {
	files, err := ioutil.ReadDir(dirPath) //读取目录下文件
	if err != nil {
		return err
	}

	for _, file := range files {
		deepFile := dirPath + GetPathDelimiter() + file.Name()
		if file.IsDir() {
			walkDirWithFilter(deepFile, allFiles, filter)
			continue
		}
		if filter(deepFile) {
			*allFiles = append(*allFiles, deepFile)
		}
	}

	return nil
}
