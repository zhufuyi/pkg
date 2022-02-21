package mconf

import (
	"reflect"
	"testing"
)

func TestDeleteYaml(t *testing.T) {
	type args struct {
		in       []byte
		selector string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "删除副本集",
			args: args{
				in:       yamlData,
				selector: ".spec.replicas",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := DeleteYaml(tt.args.in, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteYaml() error %v, wantErr %v", err, tt.wantErr)
				return
			}

			val, _ := FindYaml(gotData, tt.args.selector)
			got := Bytes2Int(val)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteYaml() got %d, want %d", got, tt.want)
			}
		})
	}
}
