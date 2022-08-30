package replacer

import (
	"embed"
	"testing"
)

//go:embed testDir
var fs embed.FS

func TestTemplate(t *testing.T) {
	//r, err := New("dir")
	//if err != nil {
	//	t.Fatal(err)
	//}
	r, err := NewWithFS("dir", fs)
	if err != nil {
		t.Fatal(err)
	}

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
	r.SetIgnoreFiles(ignoreFiles...) // 这是忽略替换文件
	r.SetReplacementFields(fields)   // 设置替换文本
	r.SetOutPath("", "test")         // 设置输出目录和名称
	err = r.SaveFiles()              // 保存替换后文件
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("save files successfully, output = %s", r.GetOutPath())
}
