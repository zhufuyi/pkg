package model

import (
	"sync"

	"github.com/zhufuyi/pkg/mysql"
	"gorm.io/gorm"
)

var (
	// ErrNotFound 空记录
	ErrNotFound = gorm.ErrRecordNotFound
)

var (
	db   *gorm.DB
	Once sync.Once
	dsn  string
)

// InitMysql 连接mysql
func InitMysql(addr string) {
	dsn = addr
	var err error
	db, err = mysql.Init(addr, mysql.WithLog())
	if err != nil {
		panic("config.Get() error: " + err.Error())
	}
}

// GetDB 返回db对象
func GetDB() *gorm.DB {
	if db == nil {
		Once.Do(func() {
			InitMysql(dsn)
		})
	}

	return db
}
