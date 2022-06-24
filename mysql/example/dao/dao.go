package dao

import (
	"github.com/zhufuyi/pkg/mysql"
)

// Dao dao
type Dao struct {
	db *mysql.DB
}

// New 实例化
func New(db *mysql.DB) *Dao {
	return &Dao{db: db}
}
