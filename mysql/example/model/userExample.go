package model

import (
	"github.com/zhufuyi/pkg/mysql"
)

// UserExample object fields mapping table
type UserExample struct {
	mysql.Model

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}

// TableName get table name
func (table *UserExample) TableName() string {
	return mysql.GetTableName(table)
}

// Create a new record
func (table *UserExample) Create(db *mysql.DB) error {
	return db.Create(table).Error
}

// Delete record
func (table *UserExample) Delete(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Delete(table).Error
}

// DeleteByID delete record by id
func (table *UserExample) DeleteByID(db *mysql.DB) error {
	return db.Where("id = ?", table.ID).Delete(table).Error
}

// Updates record
func (table *UserExample) Updates(db *mysql.DB, update mysql.KV, query interface{}, args ...interface{}) error {
	return db.Model(table).Where(query, args...).Updates(update).Error
}

// Get one record
func (table *UserExample) Get(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).First(table).Error
}

// GetByID get record by id
func (table *UserExample) GetByID(db *mysql.DB, id uint64) error {
	return db.Where("id = ?", id).First(table).Error
}

// Gets multiple records, starting from page 0
func (table *UserExample) Gets(db *mysql.DB, page *mysql.Page, query interface{}, args ...interface{}) ([]*UserExample, error) {
	out := []*UserExample{}
	err := db.Order(page.Sort()).Limit(page.Size()).Offset(page.Offset()).Where(query, args...).Find(&out).Error
	return out, err
}

// Count number of statistics
func (table *UserExample) Count(db *mysql.DB, query interface{}, args ...interface{}) (int, error) {
	count := 0
	err := db.Model(table).Where(query, args...).Count(&count).Error
	return count, err
}
