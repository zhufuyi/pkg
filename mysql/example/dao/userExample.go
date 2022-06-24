package dao

import (
	"errors"
	"strings"

	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/model"
)

// UserExample dao 对象
type UserExample struct {
	ID     uint64 `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

// CreateUserExample 创建一条记录
func (d *Dao) CreateUserExample(param *UserExample) error {
	data := &model.UserExample{
		Model:  mysql.Model{ID: param.ID},
		Name:   param.Name,
		Age:    param.Age,
		Gender: param.Gender,
	}
	return data.Create(d.db)
}

// CreateUserExamples 创建多条记录
func (d *Dao) CreateUserExamples(params []*UserExample) (int, error) {
	errStr := []string{}
	count := 0
	for _, o := range params {
		err := d.CreateUserExample(o)
		if err != nil {
			errStr = append(errStr, err.Error())
			continue
		}
		count++
	}

	if len(errStr) == 0 {
		return count, nil
	}

	return count, errors.New(strings.Join(errStr, " || "))
}

// DeleteUserExample 根据id删除一条记录
func (d *Dao) DeleteUserExample(id uint64) error {
	obj := &model.UserExample{}
	obj.ID = id
	return obj.DeleteByID(d.db)
}

// DeleteUserExamples 根据id删除多条记录
func (d *Dao) DeleteUserExamples(ids []uint64) error {
	obj := &model.UserExample{}
	return obj.Delete(d.db, "id IN (?)", ids)
}

// UpdateUserExample 更新记录
func (d *Dao) UpdateUserExample(param *UserExample) error {
	obj := &model.UserExample{}
	update := mysql.KV{}
	if param.Name != "" {
		update["name"] = param.Name
	}
	if param.Age > 0 {
		update["age"] = param.Age
	}
	if param.Gender != "" {
		update["gender"] = param.Gender
	}
	return obj.Updates(d.db, update, "id = ?", param.ID)
}

// GetUserExample 根据id获取一条记录
func (d *Dao) GetUserExample(id uint64) (*model.UserExample, error) {
	obj := &model.UserExample{}
	err := obj.GetByID(d.db, id)
	return obj, err
}

// GetUserExamplesByColumns 根据列信息筛选多条记录
// columns 列信息，列名称、列值、表达式，列之间逻辑关系
// page表示页码，从0开始, size表示每页行数, sort排序字段，默认是id倒叙，可以在字段前添加-号表示倒序，无-号表示升序
// 查询年龄大于20的男人示例：
//	columns=[]*mysql.Column{
//		{
//			Name:  "gender",
//			Value: "男",
//		},
//		{
//			Name:    "age",
//			Value:   20,
//			ExpType: mysql.Gt,
//		},
//	}
func (d *Dao) GetUserExamplesByColumns(columns []*mysql.Column, page int, pageSize int, sort string) ([]*model.UserExample, int, error) {
	query, args, err := mysql.GetQueryConditions(columns)
	if err != nil {
		return nil, 0, err
	}

	obj := &model.UserExample{}
	total, err := obj.Count(d.db, query, args...)
	if err != nil {
		return nil, total, err
	}

	pageSet := mysql.NewPage(page, pageSize, sort)
	data, err := obj.Gets(d.db, pageSet, query, args...)

	return data, total, err
}
