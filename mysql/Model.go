package mysql

import (
	"reflect"
	"time"

	"github.com/huandu/xstrings"
)

// Model 表内嵌字段
type Model struct {
	ID        uint64     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

// KV map类型
type KV = map[string]interface{}

// GetTableName 获取表名
func GetTableName(object interface{}) string {
	tableName := ""

	typeof := reflect.TypeOf(object)
	switch typeof.Kind() {
	case reflect.Ptr:
		tableName = typeof.Elem().Name()
	case reflect.Struct:
		tableName = typeof.Name()
	default:
		return tableName
	}

	return xstrings.ToSnakeCase(tableName)
}
