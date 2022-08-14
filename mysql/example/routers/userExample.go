package routers

import (
	"github.com/zhufuyi/pkg/mysql/example/handler"

	"github.com/gin-gonic/gin"
)

func userExampleRouter(group *gin.RouterGroup, h handler.UserExampleHandler) {
	group.POST("/userExample", h.Create)
	group.DELETE("/userExample/:id", h.DeleteByID)
	group.PUT("/userExample/:id", h.UpdateByID)
	group.GET("/userExample/:id", h.GetByID)
	group.GET("/userExamples", h.List)
	group.POST("/userExamples", h.List2) // 通过post查询多条记录
}
