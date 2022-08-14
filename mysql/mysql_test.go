package mysql

import (
	"fmt"
	"testing"
	"time"
)

var dsn = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInit(t *testing.T) {
	db, err := Init(dsn)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}

func TestInitNoTLS(t *testing.T) {
	db, err := Init(
		dsn,
		//WithLog(),
		WithSlowThreshold(time.Millisecond*10), // 打印超过10毫秒的日志
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
	)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, dsn=%s", err, dsn))
		return
	}

	t.Logf("%+v", db.Name())
}
