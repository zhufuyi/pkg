package mconf

import (
	"reflect"
	"testing"
)

func TestPutYaml(t *testing.T) {
	type args struct {
		in        []byte
		selector  string
		valueType string
		value     string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "修改replicas",
			args: args{
				in:        yamlData,
				selector:  ".spec.replicas",
				valueType: "int",
				value:     "3",
			},
			want:    []byte("3"),
			wantErr: false,
		},
		{
			name: "修改image",
			args: args{
				in:        yamlData,
				selector:  ".spec.template.spec.containers.(name=nginx).image",
				valueType: "string",
				value:     "nginx:1.21.3",
			},
			want:    []byte("nginx:1.21.3"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := PutYaml(tt.args.in, tt.args.selector, tt.args.valueType, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			val, _ := FindYaml(gotData, tt.args.selector)
			got := []byte(Bytes2Str(val))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutYaml() got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPutObjectYaml(t *testing.T) {
	type args struct {
		in         []byte
		selector   string
		valueTypes []string
		kvs        []string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "修改limits",
			args: args{
				in:         yamlData,
				selector:   ".spec.template.spec.containers.(name=nginx).resources.requests",
				valueTypes: []string{"string", "string"},
				kvs:        []string{"cpu=1m", "memory=10Mi"},
			},
			want:    []byte("cpu: 1mmemory: 10Mi"),
			wantErr: false,
		},
		{
			name: "修改requests",
			args: args{
				in:         yamlData,
				selector:   ".spec.template.spec.containers.(name=nginx).resources.limits",
				valueTypes: []string{"string", "string"},
				kvs:        []string{"cpu=100m", "memory=500Mi"},
			},
			want:    []byte("cpu: 100mmemory: 500Mi"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := PutObjectYaml(tt.args.in, tt.args.selector, tt.args.valueTypes, tt.args.kvs)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutObject() error %v, wantErr %v", err, tt.wantErr)
				return
			}

			val, _ := FindYaml(gotData, tt.args.selector)
			got := []byte(Bytes2Str(val))

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutObject() got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestPutDocumentYaml(t *testing.T) {
	type args struct {
		in       []byte
		selector string
		docs     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				in:       yamlData,
				selector: ".spec.template.spec.containers.(name=nginx).resources.limits",
				docs:     `{"cpu":"200m", "memory":"600Mi"}`,
			},
			want:    "cpu: 200mmemory: 600Mi",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := PutDocumentYaml(tt.args.in, tt.args.selector, tt.args.docs)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutDocumentYaml() error %v, wantErr %v", err, tt.wantErr)
				return
			}
			val, _ := FindYaml(gotData, tt.args.selector)
			got := Bytes2Str(val)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PutDocumentYaml() got %s, want %s", got, tt.want)
			}
		})
	}
}
