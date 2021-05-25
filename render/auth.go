package render

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/errcode"
	"github.com/zhufuyi/pkg/gohttp"
	"github.com/zhufuyi/pkg/logger"
)

// AuthResp 鉴权返回值
type AuthResp struct {
	User struct {
		CustomerID int64  `json:"customerId"`
		Enabled    bool   `json:"enabled"`
		Username   string `json:"username"`
		NickName   string `json:"nickName"`

		//Roles        []struct {
		//	ID        int64  `json:"id"`
		//	DataScope string `json:"dataScope"`
		//	Level     int64  `json:"level"`
		//	Name      string `json:"name"`
		//} `json:"roles"`
		//Jobs       []struct {
		//	ID   int64  `json:"id"`
		//	Name string `json:"name"`
		//} `json:"jobs"`
	} `json:"user"`
}

// 判断是否通过鉴权
func (a *AuthResp) IsAuth() bool {
	if !a.User.Enabled || a.User.CustomerID <= 0 || a.User.Username == "" {
		return false
	}

	return true
}

// CheckAuth 检查是否已鉴权
func CheckAuth(url string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		result := &AuthResp{}

		err := gohttp.GetJSON(result, url, nil, "Authorization", token)
		if err != nil || !result.IsAuth() {
			logger.Warn("check auth error", logger.String("url", url), logger.String("Authorization", token), logger.Err(err))
			Error(c, errcode.UnauthorizedTokenError)
			c.Abort()
			return
		}
		c.Set("customerID", result.User.CustomerID)
		c.Set("username", result.User.Username)
		//  处理请求
		c.Next()
	}
}

// DefaultAuth 默认鉴权id
func DefaultAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("customerID", 22)
		c.Set("username", "admin")
		//  处理请求
		c.Next()
	}
}
