package mconf

import (
	"reflect"
	"testing"
)

func TestFindYaml(t *testing.T) {
	type args struct {
		in       []byte
		selector string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "获取replicas",
			args: args{
				in:       yamlData,
				selector: ".spec.replicas",
			},
			want:    []byte("2"),
			wantErr: false,
		},
		{
			name: "获取image",
			args: args{
				in:       yamlData,
				selector: ".spec.template.spec.containers.(name=nginx).image",
			},
			want:    []byte("nginx:1.15.2"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := FindYaml(tt.args.in, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := []byte(Bytes2Str(gotData))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindYaml() got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestFindJson(t *testing.T) {
	type args struct {
		in       []byte
		selector string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "获取uid",
			args: args{
				in:       jsonData,
				selector: ".uid",
			},
			want:    []byte("liz0yRCZz"),
			wantErr: false,
		},
		{
			name: "获取options",
			args: args{
				in:       jsonData,
				selector: ".panels.(datasource=Loki).gridPos",
			},
			want:    []byte(`h: 25w: 24x: 0"y": 3`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := FindJson(tt.args.in, tt.args.selector)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := []byte(Bytes2Str(gotData))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindJson() got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestCount(t *testing.T) {
	type args struct {
		in       []byte
		selector string
		inFormat string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "获取image",
			args: args{
				in:       yamlData,
				selector: ".spec.template.spec.containers",
				inFormat: YamlFormat,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Count(tt.args.in, tt.args.selector, tt.args.inFormat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Count() got %d, want %d", got, tt.want)
			}
		})
	}
}
