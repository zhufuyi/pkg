package service

import (
	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/dao"
	"github.com/zhufuyi/pkg/mysql/example/model"
)

// CreateUserExampleRequest 请求参数
type CreateUserExampleRequest struct {
	Name   string `form:"name" binding:"min=1"`
	Age    int    `form:"age" binding:"gt=0,lt=120"`
	Gender string `form:"gender" binding:"min=1"`
}

// CreateUserExample 创建
func (s *Service) CreateUserExample(req *CreateUserExampleRequest) error {
	data := &dao.UserExample{
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	}
	return s.dao.CreateUserExample(data)
}

// DeleteUserExampleRequest // 删除一个id时，从url参数
type DeleteUserExampleRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// DeleteUserExamplesRequest 删除多个id时，从body获取
type DeleteUserExamplesRequest struct {
	IDs []uint64 `form:"ids" binding:"min=1"`
}

// DeleteUserExample 删除记录
func (s *Service) DeleteUserExample(ids ...uint64) error {
	if len(ids) == 1 {
		return s.dao.DeleteUserExample(ids[0])
	}

	// 批量删除
	var (
		total = len(ids)
		start = 0
		end   = 0
		size  = 100
	)
	for start < total {
		end += size
		if end > total {
			end = total
		}
		err := s.dao.DeleteUserExamples(ids[start:end])
		if err != nil {
			return err
		}
		start = end
	}

	return nil
}

// UpdateUserExampleRequest 请求参数
type UpdateUserExampleRequest struct {
	ID     uint64 `form:"id" binding:"gt=0"`
	Name   string `form:"name" binding:""`
	Age    int    `form:"age" binding:""`
	Gender string `form:"gender" binding:""`
}

// UpdateUserExample 更新
func (s *Service) UpdateUserExample(req *UpdateUserExampleRequest) error {
	return s.dao.UpdateUserExample(&dao.UserExample{
		ID:     req.ID,
		Name:   req.Name,
		Age:    req.Age,
		Gender: req.Gender,
	})
}

// GetUserExampleRequest 请求参数
type GetUserExampleRequest struct {
	ID uint64 `form:"id" binding:"gt=0"`
}

// GetUserExample 根据id获取一条记录
func (s *Service) GetUserExample(req *GetUserExampleRequest) (*model.UserExample, error) {
	return s.dao.GetUserExample(req.ID)
}

// GetUserExamplesRequest 请求参数
type GetUserExamplesRequest struct {
	// get url原生请求参数，用来填充Exps、Logics的默认值
	// 如果ParamSrc为空，必须满足len(Keys)=len(Values)=len(Exps)=len(Logics)
	ParamSrc string `form:"-" binding:"-"`

	Keys   []string      `form:"k" binding:"-"`
	Values []interface{} `form:"v" binding:"-"`
	Exps   []string      `form:"exp" binding:"-"`
	Logics []string      `form:"logic" binding:"-"`

	Page int    `form:"page" binding:"gte=0"`
	Size int    `form:"size" binding:"gt=0"`
	Sort string `form:"sort" binding:"-"`
}

// GetUserExamples 获取多条记录
func (s *Service) GetUserExamples(req *GetUserExamplesRequest) ([]*model.UserExample, int, error) {
	columns, err := mysql.GetColumns(req.Keys, req.Values, req.Exps, req.Logics, req.ParamSrc)
	if err != nil {
		return nil, 0, err
	}

	return s.dao.GetUserExamplesByColumns(columns, req.Page, req.Size, req.Sort)
}

// ------------------------------------------------------------------------------------------

// 通过post提交表单查询

type column struct {
	Name  string      `json:"name"`  // 列名
	Value interface{} `json:"value"` // 值
	Exp   string      `json:"exp"`   // 表达式，值为空时默认为eq，有eq、neq、gt、gte、lt、lte、like七种类型
	Logic string      `json:"logic"` // 逻辑类型，值为空时默认为and，有and、or两种类型
}

// GetUserExamplesRequest2 请求参数
type GetUserExamplesRequest2 struct {
	Columns []column `json:"columns"`

	Page int    `form:"page" binding:"gte=0" json:"page"`
	Size int    `form:"size" binding:"gt=0" json:"size"`
	Sort string `form:"sort" binding:"" json:"sort"`
}

// GetUserExamples2 获取多条记录
func (s *Service) GetUserExamples2(req *GetUserExamplesRequest2) ([]*model.UserExample, int, error) {
	var columns []*mysql.Column
	for _, v := range req.Columns {
		if v.Value == "" {
			continue
		}
		columns = append(columns, &mysql.Column{
			Name:      v.Name,
			Value:     v.Value,
			ExpType:   v.Exp,
			LogicType: v.Logic,
		})
	}

	return s.dao.GetUserExamplesByColumns(columns, req.Page, req.Size, req.Sort)
}
