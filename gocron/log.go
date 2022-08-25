package gocron

import (
	"github.com/zhufuyi/pkg/logger"

	"go.uber.org/zap"
)

type myLog struct {
	zapLog *zap.Logger
}

func (l *myLog) Info(msg string, keysAndValues ...interface{}) {
	if msg == "wake" { // 忽略wake
		return
	}
	msg = "cron_" + msg
	fields := parseKVs(keysAndValues)
	if l.zapLog != nil {
		l.zapLog.Info(msg, fields...)
	} else {
		logger.Info(msg, fields...)
	}
}
func (l *myLog) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := parseKVs(keysAndValues)
	fields = append(fields, zap.String("err", err.Error()))
	msg = "cron_" + msg
	if l.zapLog != nil {
		l.zapLog.Error(msg, fields...)
	} else {
		logger.Error(msg, fields...)
	}
}

func parseKVs(kvs interface{}) []zap.Field {
	var fields []zap.Field

	infos, ok := kvs.([]interface{})
	if !ok {
		return fields
	}

	l := len(infos)
	if l%2 == 1 {
		return fields
	}

	for i := 0; i < l; i += 2 {
		key := infos[i].(string) //nolint
		value := infos[i+1]

		// 把id替换为任务名称
		if key == "entry" {
			if id, ok := value.(cron.EntryID); ok {
				key = "task"
				if v, isExist := idName.Load(id); isExist {
					value = v
				}
			}
		}

		fields = append(fields, zap.Any(key, value))
	}

	return fields
}
