package service

import (
	"context"

	"github.com/zhufuyi/pkg/mysql/example/dao"
)

// Service 方法
type Service struct {
	ctx context.Context
	dao *dao.Dao
}

// New 实例化
func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	//svc.dao = dao.New(global.MysqlDB)
	return svc
}
