package dao

import (
	"fmt"
	"testing"

	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
)

var dao *Dao
var addr = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	db, err := mysql.Init(addr, mysql.WithLog(logger.Get()))
	if err != nil {
		panic(fmt.Sprintf("connect to mysql failed, err=%s, addr=%s", err, addr))
	}

	dao = &Dao{db: db}
}

func TestDao_CreateUserExample(t *testing.T) {
	err := dao.CreateUserExample(&UserExample{
		Name:   "黄忠",
		Age:    22,
		Gender: "男",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDao_CreateUserExamples(t *testing.T) {
	users := []*UserExample{{
		Name:   "关平1",
		Age:    15,
		Gender: "男",
	},
		{
			Name:   "关平2",
			Age:    16,
			Gender: "男",
		},
		{
			Name:   "关平2",
			Age:    16,
			Gender: "男",
		},
	}
	n, err := dao.CreateUserExamples(users)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log("create success", n)
}

func TestDao_UpdateUserExample(t *testing.T) {
	err := dao.UpdateUserExample(&UserExample{
		ID:  16,
		Age: 23,
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDao_GetUserExample(t *testing.T) {
	user, err := dao.GetUserExample(1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", user)
}

func TestDao_DeleteUserExample(t *testing.T) {
	err := dao.DeleteUserExample(17)
	if err != nil {
		t.Error(err)
	}
}

func TestDao_GetUserExamplesByColumns(t *testing.T) {
	type fields struct {
		db *mysql.DB
	}
	type args struct {
		columns  []*mysql.Column
		page     int
		pageSize int
		sort     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "1 column",
			fields: fields{dao.db},
			args: args{
				columns: []*mysql.Column{
					{
						Name:      "name",
						Value:     "刘备",
						ExpType:   "",
						LogicType: "",
					},
				},
				page:     0,
				pageSize: 10,
				sort:     "",
			},
		},
		{
			name:   "1 column like",
			fields: fields{dao.db},
			args: args{
				columns: []*mysql.Column{
					{
						Name:    "name",
						Value:   "乔",
						ExpType: "like",
					},
				},
				page:     0,
				pageSize: 10,
				sort:     "",
			},
		},
		{
			name:   "2 column and",
			fields: fields{dao.db},
			args: args{
				columns: []*mysql.Column{
					{
						Name:  "gender",
						Value: "女",
					},
					{
						Name:    "age",
						Value:   20,
						ExpType: mysql.Lt,
					},
				},
				page:     0,
				pageSize: 10,
				sort:     "",
			},
		},
		{
			name:   "2 column or",
			fields: fields{dao.db},
			args: args{
				columns: []*mysql.Column{
					{
						Name:  "name",
						Value: "刘备",
					},
					{
						Name:      "age",
						Value:     20,
						ExpType:   mysql.Lt,
						LogicType: mysql.OR,
					},
				},
				page:     0,
				pageSize: 10,
				sort:     "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dao{
				db: tt.fields.db,
			}
			got, _, err := d.GetUserExamplesByColumns(tt.args.columns, tt.args.page, tt.args.pageSize, tt.args.sort)
			if err != nil {
				t.Error(err)
				return
			}
			for _, user := range got {
				t.Logf("%+v", user)
			}
		})
	}
}
