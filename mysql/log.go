package mysql

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
)

type gormLogger struct {
	Log *zap.Logger
}

func newGormLogger(log *zap.Logger) *gormLogger {
	return &gormLogger{log}
}

// Print 打印日志，注：gorm在更新时候如果更新字段为空，不会调用Print打印日志
func (g *gormLogger) Print(values ...interface{}) {
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

		g.Log.WithOptions(zap.AddCallerSkip(skip)).Info("gorm",
			zap.Int64("ns", int64(values[2].(time.Duration))),
			zap.Int64("rows", values[5].(int64)),
			zap.String("sql", sqlStr),
			zap.String("values", paramStr),
		)
	}
}
