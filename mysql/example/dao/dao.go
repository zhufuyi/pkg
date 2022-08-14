package dao

import "gorm.io/gorm"

// Dao 对象
type Dao struct {
	db *gorm.DB
}

// NewDao 新建dao示例
func NewDao(db *gorm.DB) *Dao {
	return &Dao{
		db: db,
	}
}
