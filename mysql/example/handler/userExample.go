package handler

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/zhufuyi/pkg/gin/response"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/dao"
	"github.com/zhufuyi/pkg/mysql/example/errcode"
	"github.com/zhufuyi/pkg/mysql/example/model"
)

var _ UserExampleHandler = (*userExampleHandler)(nil)

// UserExampleHandler 定义handler接口
type UserExampleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)
	List2(c *gin.Context)
}

type userExampleHandler struct {
	iDao dao.UserExampleDao
}

// NewUserExampleHandler 创建handler接口
func NewUserExampleHandler(iDao dao.UserExampleDao) UserExampleHandler {
	return &userExampleHandler{
		iDao: iDao,
	}
}

// Create 创建
func (h *userExampleHandler) Create(c *gin.Context) {
	form := &CreateUserExampleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err))
		response.Error(c, errcode.InvalidParams)
		return
	}

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.InternalServerError)
		return
	}

	err = h.iDao.Create(c, userExample)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.ErrCreateUserExample)
		return
	}

	response.Success(c, gin.H{"id": userExample.ID})
}

// DeleteByID 根据id删除一条记录
func (h *userExampleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	err := h.iDao.DeleteByID(c, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id))
		response.Error(c, errcode.ErrDeleteUserExample)
		return
	}

	response.Success(c)
}

// UpdateByID 根据id更新信息
func (h *userExampleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	form := &UpdateUserExampleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err))
		response.Error(c, errcode.InvalidParams)
		return
	}
	form.ID = id

	userExample := &model.UserExample{}
	err = copier.Copy(userExample, form)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.InternalServerError)
		return
	}

	err = h.iDao.UpdateByID(c, userExample)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.ErrUpdateUserExample)
		return
	}

	response.Success(c)
}

// GetByID 根据id获取一条记录
func (h *userExampleHandler) GetByID(c *gin.Context) {
	idstr, id, isAbout := getIDFromPath(c)
	if isAbout {
		return
	}

	userExample, err := h.iDao.GetByID(c, id)
	if err != nil {
		if err.Error() == mysql.ErrNotFound.Error() {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id))
			response.Error(c, errcode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id))
			response.Error(c, errcode.ErrGetUserExample)
		}
		return
	}

	data := &GetUserExampleByIDRespond{}
	err = copier.Copy(data, userExample)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("id", id))
		response.Error(c, errcode.InternalServerError)
		return
	}
	data.ID = idstr

	response.Success(c, gin.H{"userExample": data})
}

// List 获取多条记录
// 通过url参数作为查询条件，支持任意多个字段，下面以user表为例子get请求参数，不同条件查询第0页20条记录，默认是id降序
// 没有条件查询 ?page=0&size=20
// 名称查询 ?page=0&size=20&k=name&v=张三
// 名称模糊查询 ?page=0&size=20&k=name&v=张&exp=like
// 年龄为18岁的男性 ?page=0&size=20&k=age&v=22&gender=1
// 年龄小于18或者大于60的人 ?page=0&size=20&k=age&v=18&exp=gt&logic=or&k=age&v=60&exp=lt
// 排序可以在参数添加sort字段，例如sort=id表示id升序，sort=-id表示id降序
func (h *userExampleHandler) List(c *gin.Context) {
	form := &GetUserExamplesRequest{}
	err := c.ShouldBindQuery(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err))
		response.Error(c, errcode.InvalidParams)
		return
	}
	form.URLParams = c.Request.URL.RawQuery

	var values []interface{}
	for _, v := range form.Values {
		values = append(values, v)
	}
	columns, err := mysql.GetColumns(form.Keys, values, form.Exps, form.Logics, form.URLParams)
	if err != nil {
		logger.Warn("GetColumns error: ", logger.Err(err))
		response.Error(c, errcode.InvalidParams)
		return
	}

	userExamples, total, err := h.iDao.GetByColumns(c, columns, form.Page, form.Size, form.Sort)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.ErrGetUserExample)
		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.InternalServerError)
		return
	}

	response.Success(c, gin.H{
		"userExamples": data,
		"total":        total,
	})
}

// List2 通过post获取多条记录
func (h *userExampleHandler) List2(c *gin.Context) {
	form := &GetUserExamplesRequest2{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err))
		response.Error(c, errcode.InvalidParams)
		return
	}

	var columns []*mysql.Column
	for _, v := range form.Columns {
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

	userExamples, total, err := h.iDao.GetByColumns(c, columns, form.Page, form.Size, form.Sort)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.ErrGetUserExample)
		return
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		logger.Error("Copy error", logger.Err(err), logger.Any("form", form))
		response.Error(c, errcode.InternalServerError)
		return
	}

	response.Success(c, gin.H{
		"userExamples": data,
		"total":        total,
	})
}

// ------------------------------------定义请求参数和返回结果------------------------------

// CreateUserExampleRequest 创建请求参数
type CreateUserExampleRequest struct {
	// binding使用说明 https://github.com/go-playground/validator

	Name     string `json:"name" binding:"min=2"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password" binding:"md5"`
	Phone    string `form:"phone" binding:"e164"`
	Age      int    `form:"age" binding:"gt=0,lt=120"`
	Gender   int    `form:"gender" binding:"gte=0,lte=2"`
}

// UpdateUserExampleByIDRequest 更新请求参数
type UpdateUserExampleByIDRequest struct {
	ID uint64 `json:"id" binding:"-"`

	CreateUserExampleRequest
}

type simpleUpdateUserExample struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Age       int       `json:"age"`
	Gender    int       `json:"gender"`
	Status    int       `json:"status"`
	LoginAt   int64     `json:"login_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUserExampleByIDRespond 返回数据
type GetUserExampleByIDRespond struct {
	ID string `json:"id"`

	simpleUpdateUserExample
}

// ListUserExamplesRespond 返回列表数据
type ListUserExamplesRespond []struct {
	GetUserExampleByIDRespond
}

// GetUserExamplesRequest 请求参数
type GetUserExamplesRequest struct {
	Page int    `form:"page" binding:"gte=0"`
	Size int    `form:"size" binding:"gt=0"`
	Sort string `form:"sort" binding:"-"`

	// 参数填写方式一：从request请求url中获取参数(form.URLParams = c.Request.URL.RawQuery)，
	// 用来自动填充exp、logic的默认值，为了在url参数减少填写exp和logic的默认值，例如url参数?page=0&size=20&k=age&exp=gt&v=22&k=gender&v=1，表示查询年龄大于22岁的男性
	// 参数填写方式二：没有从请求url中获取参数，也就是ParamSrc为空时，请求url参数必须满足len(k)=len(v)=len(exp)=len(logic)，
	// 可以同时存在多个，也可以同时不存在，例如url参数?page=0&size=20&k=age&v=22&exp=gt&logic=and&k=gender&v=1&exp=eq&logic=and，也是表示查询年龄大于22岁的男性
	// 两种url参数都是合法，建议使用第一种方式
	URLParams string   `form:"-" binding:"-"`
	Keys      []string `form:"k" binding:"-"`
	Values    []string `form:"v" binding:"-"`
	Exps      []string `form:"exp" binding:"-"`
	Logics    []string `form:"logic" binding:"-"`
}

// 通过post方法提交表单进行查询
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

// --------------------------------除了handler的其他功能函数------------------------------

func getIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		logger.Error("ParseUint error: ", logger.String("idStr", idStr))
		response.Error(c, errcode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertUserExamples(fromValues []model.UserExample) ([]GetUserExampleByIDRespond, error) {
	var toValues = []GetUserExampleByIDRespond{}
	for _, v := range fromValues {
		data := GetUserExampleByIDRespond{}
		err := copier.Copy(&data, &v)
		if err != nil {
			return nil, err
		}
		data.ID = strconv.FormatUint(v.ID, 10)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
