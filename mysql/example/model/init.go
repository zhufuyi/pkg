package model

import (
	"sync"

	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/config"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	Once sync.Once
)

// InitMysql 连接mysql
func InitMysql() {
	var err error
	db, err = mysql.Init(config.Get().MysqlURL, mysql.WithLog())
	if err != nil {
		panic("config.Get() error: " + err.Error())
	}
}

// GetDB 返回db对象
func GetDB() *gorm.DB {
	if db == nil {
		Once.Do(func() {
			InitMysql()
		})
	}

	return db
}
