package render

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// BindJSON 解析参数
func BindJSON(c *gin.Context, obj interface{}) (err error) {
	return binding.JSON.Bind(c.Request, obj)
}

// BindForm 解析参数
func BindForm(c *gin.Context, obj interface{}) (err error) {
	return binding.Form.Bind(c.Request, obj)
}

// QueryInt 获取URL参数
func QueryInt(c *gin.Context, key string) (int, error) {
	s, ok := c.GetQuery(key)
	if ok {
		return strconv.Atoi(s)
	}
	return 0, errors.New("empty key")
}

// QueryAll 获取URL所有参数
func QueryAll(c *gin.Context) map[string]interface{} {
	vals := c.Request.URL.Query()
	data := make(map[string]interface{})

	for k, v := range vals {
		data[k] = v[0]
	}

	return data
}

// QueryInt64 获取URL参数
func QueryInt64(c *gin.Context, key string) (int64, error) {
	s, ok := c.GetQuery(key)
	if ok {
		return strconv.ParseInt(s, 10, 64)
	}

	return 0, errors.New("empty key")
}

// ParamInt64 获取Param参数
func ParamInt64(c *gin.Context, key string) (int64, error) {
	s := c.Param(key)
	if s != "" {
		return strconv.ParseInt(s, 10, 64)
	}

	return 0, errors.New("empty key")
}
