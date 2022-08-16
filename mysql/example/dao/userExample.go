package dao

import (
	"context"
	"errors"

	"github.com/zhufuyi/pkg/mysql/example/model"
	"github.com/zhufuyi/pkg/mysql/query"
)

var _ UserExampleDao = (*userExampleDao)(nil)

// UserExampleDao 定义dao接口
type UserExampleDao interface {
	Create(ctx context.Context, table *model.UserExample) error
	DeleteByID(ctx context.Context, id uint64) error
	UpdateByID(ctx context.Context, table *model.UserExample) error
	GetByID(ctx context.Context, id uint64) (*model.UserExample, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]model.UserExample, int64, error)
}

type userExampleDao struct {
	*Dao
}

// NewUserExampleDao 创建dao接口
func NewUserExampleDao(dao *Dao) UserExampleDao {
	return &userExampleDao{dao}
}

// Create 创建一条记录，插入记录后，id值被回写到table中
func (d *userExampleDao) Create(ctx context.Context, table *model.UserExample) error {
	return d.db.WithContext(ctx).Create(table).Error
}

// DeleteByID 根据id删除一条记录
func (d *userExampleDao) DeleteByID(ctx context.Context, id uint64) error {
	return d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserExample{}).Error
}

// Deletes 根据id删除多条记录
func (d *userExampleDao) Deletes(ctx context.Context, ids []uint64) error {
	return d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.UserExample{}).Error
}

// UpdateByID 根据id更新记录
func (d *userExampleDao) UpdateByID(ctx context.Context, table *model.UserExample) error {
	if table.ID < 1 {
		return errors.New("id cannot be less than 0")
	}

	update := map[string]interface{}{}
	if table.Name != "" {
		update["name"] = table.Name
	}
	if table.Password != "" {
		update["password"] = table.Password
	}
	if table.Email != "" {
		update["email"] = table.Email
	}
	if table.Phone != "" {
		update["phone"] = table.Phone
	}
	if table.Avatar != "" {
		update["avatar"] = table.Avatar
	}
	if table.Age > 0 {
		update["age"] = table.Age
	}
	if table.Gender > 0 {
		update["gender"] = table.Gender
	}
	if table.LoginAt > 0 {
		update["login_at"] = table.LoginAt
	}

	return d.db.WithContext(ctx).Model(table).Where("id = ?", table.ID).Updates(update).Error
}

// GetByID 根据id获取一条记录
func (d *userExampleDao) GetByID(ctx context.Context, id uint64) (*model.UserExample, error) {
	table := &model.UserExample{}

	err := d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
	if err != nil {
		return nil, err
	}
	return table, nil
}

// GetByColumns 根据列信息筛选多条记录
// columns 列信息，列名称、列值、表达式，列之间逻辑关系
// page表示页码，从0开始, size表示每页行数, sort排序字段，默认是id倒叙，可以在字段前添加-号表示倒序，无-号表示升序
// 示例：查询年龄大于20的男性
//
//	columns=[]*mysql.Column{
//		{
//			Name:  "gender",
//			Value: "男",
//		},
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//	}
func (d *userExampleDao) GetByColumns(ctx context.Context, params *query.Params) ([]model.UserExample, int64, error) {
	query, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = d.db.WithContext(ctx).Model(&model.UserExample{}).Where(query, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, total, nil
	}

	tables := []model.UserExample{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(query, args...).Find(&tables).Error
	if err != nil {
		return nil, 0, err
	}

	return tables, total, err
}
