package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/gin/errcode"
	"github.com/zhufuyi/pkg/gin/render"
	"github.com/zhufuyi/pkg/jwt"
	"github.com/zhufuyi/pkg/logger"
)

// JWT 中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authentication := c.GetHeader("Authentication")
		if authentication == "" || len(authentication) < 20 {
			logger.Error("token is empty")
			render.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}
		token := authentication[7:]
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			logger.Error("VerifyToken error", logger.Err(err))
			render.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}

		uid := c.GetHeader("X-Uid")
		if uid == "" {
			logger.Error("uid is empty")
			render.Error(c, errcode.InvalidParams)
			c.Abort()
			return
		}

		// 判断你是否uid一致
		if claims.Uid != uid {
			logger.Error("can't confirm it's you", logger.String("claims.uid", claims.Uid), logger.String("header.uid", uid))
			render.Error(c, errcode.Unauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
