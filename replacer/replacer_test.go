package replacer

import (
	"embed"
	"testing"
)

//go:embed testDir
var fs embed.FS

func TestNewWithFS(t *testing.T) {
	type args struct {
		newFun func() Replacer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "new",
			args: args{
				newFun: func() Replacer {
					replacer, err := New("testDir")
					if err != nil {
						panic(err)
					}
					return replacer
				},
			},
			wantErr: false,
		},

		{
			name: "new fs",
			args: args{
				newFun: func() Replacer {
					replacer, err := NewWithFS("testDir", fs)
					if err != nil {
						panic(err)
					}
					return replacer
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.args.newFun()

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
			err := r.SaveFiles()             // 保存替换后文件
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWithFS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("save files successfully, output = %s", r.GetOutPath())
		})
	}
}