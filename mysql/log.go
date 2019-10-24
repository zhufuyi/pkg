package mysql

import (
	"database/sql/driver"
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/zhufuyi/logger"
)

// gormLogger gorm日志
type gormLogger struct {
	Log *logger.ZapLogger
}

func newGormLogger() *gormLogger {
	return &gormLogger{logger.WithFields()}
}

// Format Log
var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

func (l gormLogger) Print(values ...interface{}) {
	if len(values) > 2 {
		skip := 7
		level := values[0]
		source := fmt.Sprintf("%s", values[1])
		params := strings.Split(source, ":")
		source = fmt.Sprintf("%s:%s", path.Base(params[0]), params[1])
		messages := []logger.Field{}

		if level == "sql" {
			var formatedValues []interface{}
			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if _, ok := value.([]byte); ok {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", "[binary]"))
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
						} else {
							formatedValues = append(formatedValues, "NULL")
						}
					} else {
						formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formatedValues = append(formatedValues, fmt.Sprintf("'%v'", value))
				}
			}

			messages = append(messages,
				logger.Int64("timeUS", values[2].(time.Duration).Nanoseconds()/1000),
				logger.Int64("rows", values[5].(int64)), // 返回或有效影响的行数
				logger.String("sql", fmt.Sprintf(sqlRegexp.ReplaceAllString(values[3].(string), "%v"), formatedValues...)),
			)

			// INSERT和SELECT对应的skip=7，UPDATE和DELETE对应skip=8
			if strings.Contains(values[3].(string), "UPDATE") || strings.Contains(values[3].(string), "DELETE") {
				skip = 8
			}
		} else {
			skip = 8
			messages = append(messages, logger.Any("sql", values[2:]))
		}

		logger.GetLogger(skip).Debug("mysql info", messages...)
	}
}
