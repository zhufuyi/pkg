package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/zhufuyi/pkg/logger"
)

// gormLogger gorm日志
type gormLogger struct {
	Log *logger.ZapLogger
}

func newGormLogger() *gormLogger {
	return &gormLogger{logger.WithFields()}
}

func (l *gormLogger) Print(values ...interface{}) {
	if len(values) > 5 {
		skip := 8
		sqlStr := values[3].(string)
		if len(sqlStr) > 10 {
			if strings.Contains(sqlStr[:10], "DELETE") || strings.Contains(sqlStr[:10], "UPDATE") {
				skip = 9
			}
		}

		paramStr := fmt.Sprintf("%v", values[4])
		if len(paramStr) > 300 {
			paramStr = paramStr[0:300] + " ......"
		}

		logger.GetLogger(skip).Info("gorm",
			logger.Int64("ns", int64(values[2].(time.Duration))),
			logger.Int64("rows", values[5].(int64)),
			logger.String("sql", sqlStr),
			logger.String("values", paramStr),
		)
	}
}
