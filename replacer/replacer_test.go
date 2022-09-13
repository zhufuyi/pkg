package replacer

import (
	"embed"
	"testing"
)

//go:embed testDir
var fs embed.FS

func TestTemplate(t *testing.T) {
	//r, err := New("testDir")
	//if err != nil {
	//	t.Fatal(err)
	//}
	r, err := NewWithFS("testDir", fs)
	if err != nil {
		t.Fatal(err)
	}

	subDirs := []string{}
	ignoreDirs := []string{}
	ignoreFiles := []string{}
	fields := []Field{
		{
			Old: "1234",
			New: "....",
		},
		{
			Old:             "abcdef",
			New:             "hello_",
			IsCaseSensitive: true,
		},
	}
	r.SetSubDirs(subDirs...)         // 只处理指定子目录，为空时表示指定全部文件
	r.SetIgnoreFiles(ignoreDirs...)  // 忽略替换目录
	r.SetIgnoreFiles(ignoreFiles...) // 忽略替换文件
	r.SetReplacementFields(fields)   // 设置替换文本
	r.SetOutDir("", "test")          // 设置输出目录和名称
	err = r.SaveFiles()              // 保存替换后文件
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("save files successfully, output = %s", r.GetOutPath())
}
