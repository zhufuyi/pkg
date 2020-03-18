package gofile

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

func TestGetRunPath(t *testing.T) {
	fmt.Println(GetRunPath())
}

func TestListFiles(t *testing.T) {
	//files, err := ListFiles(GetRunPath())
	//files, err := ListFiles("C:\\Work\\Golang\\Project\\src\\ts")
	files, err := ListFiles("C:\\Work\\Golang\\Package\\src")
	//files, err := ListFiles(".")
	if err != nil {
		t.Error(err)
		return
	}

	for _, file := range files {
		fmt.Println(file)
	}
	fmt.Println(len(files))
}

func TestListFilesWithFilter(t *testing.T) {
	//files, err := ListFilesWithFilter(GetRunPath(), MatchContain(".exe"))
	//files, err := ListFilesWithFilter("C:\\Work\\Golang\\Project\\src\\ts", MatchSuffix("test.go"))
	//files, err := ListFilesWithFilter("C:\\Work\\Golang\\Package\\src", MatchSuffix("test.go"))
	files, err := ListFilesWithFilter("..", MatchSuffix(".go"))
	if err != nil {
		t.Error(err)
		return
	}

	for _, file := range files {
		fmt.Println(file)
	}

	fmt.Println(len(files))
}

func TestListDirsAndFiles(t *testing.T) {
	df, err := ListDirsAndFiles("..")
	if err != nil {
		t.Error(err)
		return
	}
	pp.Println(df)
}

func TestDeleteDir(t *testing.T) {
	vals, err := DeleteDir("C:\\Users\\zhuya\\AppData\\Local\\Temp\\strategies\\dirs")
	if err != nil {
		t.Error(err)
	}
	pp.Println(vals)
}
