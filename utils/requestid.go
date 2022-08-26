package utils

import (
	"github.com/zhufuyi/pkg/gin/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestID zap logger request ID
func RequestID(c *gin.Context) zap.Field {
	return zap.String(middleware.ContextRequestIDKey, middleware.GetRequestIDFromContext(c))
}
