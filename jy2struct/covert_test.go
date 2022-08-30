package jy2struct

import "testing"

func TestCovert(t *testing.T) {
	type args struct {
		args *Args
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "json to struct from data",
			args: args{args: &Args{
				Data:   `{"name":"zhangsan","age":22}`,
				Format: "json",
			}},
			wantErr: false,
		},
		{
			name: "yaml to struct from data",
			args: args{args: &Args{
				Data: `
name": zhangsan
first_name: zhu
age: 10
`,
				Format: "yaml",
			}},
			wantErr: false,
		},
		{
			name: "json to struct from file",
			args: args{args: &Args{
				InputFile: "test.json",
				Format:    "json",
				SubStruct: true,
			}},
			wantErr: false,
		},
		{
			name: "yaml to struct from file",
			args: args{args: &Args{
				InputFile: "test.yaml",
				Format:    "yaml",
				SubStruct: true,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Covert(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Covert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
		})
	}
}
