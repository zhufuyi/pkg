package email

import (
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "qq",
			args: args{
				username: "xxxxxx@qq.com",
				password: "xxx",
			},
			want:    "smtp.qq.com",
			wantErr: false,
		},
		{
			name: "126",
			args: args{
				username: "xxxxxx@126.com",
				password: "xxx",
			},
			want:    "smtp.126.com",
			wantErr: false,
		},
		{
			name: "163",
			args: args{
				username: "xxxxxx@163.com",
				password: "",
			},
			want:    "smtp.163.com",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Init(tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Host, tt.want) {
				t.Errorf("Init() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailClient_SendMessage(t *testing.T) {
	// 发送内容
	msg := &Message{
		To:          []string{"xxxxxx@qq.com"},
		Cc:          nil,
		Subject:     "title-demo",
		ContentType: "text/plain",
		Content:     "邮件内容demo-01",
		Attach:      "",
	}

	// qq邮箱发送
	client, _ := Init("xxxxxx@qq.com", "xxxxxx")
	err := client.SendMessage(msg)
	if err != nil {
		t.Log(err)
	}

	// 126邮箱发送
	client, _ = Init("xxxxxx@126.com", "xxxxxx")
	err = client.SendMessage(msg)
	if err != nil {
		t.Log(err)
	}

	// 163邮箱发送
	client, _ = Init("xxxxxx@63.com", "xxxxxx")
	err = client.SendMessage(msg)
	if err != nil {
		t.Log(err)
	}
}
