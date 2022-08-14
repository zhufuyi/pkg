package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zhufuyi/pkg/gin/middleware"
	"github.com/zhufuyi/pkg/gin/validator"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/mysql/example/config"
	"github.com/zhufuyi/pkg/mysql/example/dao"
	"github.com/zhufuyi/pkg/mysql/example/handler"
)

// NewRouter 实例化路由
func NewRouter(idao *dao.Dao) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))
	r.Use(middleware.Cors())
	binding.Validator = validator.Init()

	if config.Get().IsEnableProfile {
		pprof.Register(r)
	}

	apiv1 := r.Group("/api/v1")

	// 注册路由
	userExampleRouter(apiv1, handler.NewUserExampleHandler(dao.NewUserExampleDao(idao)))

	return r
}
