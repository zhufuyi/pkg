package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/dao"
	"github.com/zhufuyi/pkg/mysql/example/model"
)

var daoDb *dao.Dao
var addr = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	db, err := mysql.Init(addr, mysql.WithLog())
	if err != nil {
		panic(fmt.Sprintf("connect to mysql failed, err=%s, addr=%s", err, addr))
	}

	daoDb = dao.New(db)
}

func TestService_GetUserExamples(t *testing.T) {
	type fields struct {
		ctx context.Context
		dao *dao.Dao
	}
	type args struct {
		req *GetUserExamplesRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.UserExample
		want1   int
		wantErr bool
	}{
		{
			name:   "0 column",
			fields: fields{context.Background(), daoDb},
			args: args{
				req: &GetUserExamplesRequest{
					ParamSrc: "page=0&size=10",
					Keys:     nil,
					Exps:     nil,
					Values:   nil,
					Logics:   nil,
					Page:     0,
					Size:     10,
					Sort:     "",
				},
			},
			want:    nil,
			want1:   9,
			wantErr: false,
		},
		{
			name:   "1 column gt",
			fields: fields{context.Background(), daoDb},
			args: args{
				req: &GetUserExamplesRequest{
					ParamSrc: "k=age&exp=gt&v=20&page=0&size=10",
					Keys:     []string{"age"},
					Exps:     nil,
					Values:   []interface{}{"20"},
					Logics:   nil,
					Page:     0,
					Size:     10,
					Sort:     "",
				},
			},
			want:    nil,
			want1:   9,
			wantErr: false,
		},
		{
			name:   "3 column neq gt or",
			fields: fields{context.Background(), daoDb},
			args: args{
				req: &GetUserExamplesRequest{
					ParamSrc: "k=name&v=刘备&exp=neq&k=age&v=20&exp=gt&k=gender&v=男&logic=or&page=0&size=10",
					Keys:     []string{"name", "age", "gender"},
					Exps:     nil,
					Values:   []interface{}{"刘备", "20", "男"},
					Logics:   nil,
					Page:     0,
					Size:     10,
					Sort:     "",
				},
			},
			want:    nil,
			want1:   9,
			wantErr: false,
		},
		{
			name:   "paramSrc is empty",
			fields: fields{context.Background(), daoDb},
			args: args{
				req: &GetUserExamplesRequest{
					ParamSrc: "",
					Keys:     []string{"age"},
					Exps:     []string{"gt"},
					Values:   []interface{}{"20"},
					Logics:   []string{""},
					Page:     0,
					Size:     10,
					Sort:     "",
				},
			},
			want:    nil,
			want1:   9,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				ctx: tt.fields.ctx,
				dao: tt.fields.dao,
			}
			got, _, err := s.GetUserExamples(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserExamples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, user := range got {
				t.Logf("%+v", user)
			}
		})
	}
}
