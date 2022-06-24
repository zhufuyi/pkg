package mysql

import (
	"errors"
	"fmt"
	"strings"
)

// Exp 表达式
type Exp = string

const (
	Eq   Exp = "eq"
	Neq  Exp = "neq"
	Gt   Exp = "gt"
	Gte  Exp = "gte"
	Lt   Exp = "lt"
	Lte  Exp = "lte"
	Like Exp = "like"
)

var ExpMap = map[Exp]string{
	Eq:   " = ",
	Neq:  " <> ",
	Gt:   " > ",
	Gte:  " >= ",
	Lt:   " < ",
	Lte:  " <= ",
	Like: " LIKE ",
}

// Logic 逻辑类型
type Logic = string

const (
	AND Logic = "and"
	OR  Logic = "or"
)

var logicMap = map[Logic]string{
	AND: " AND ",
	OR:  " OR ",
}

// Column 表的列查询信息
type Column struct {
	Name      string      // 列名
	Value     interface{} // 值
	ExpType   Exp         // 表达式，值为空时默认为eq，有eq、neq、gt、gte、lt、lte、like七种类型
	LogicType Logic       // 逻辑类型，值为空时默认为and，有and、or两种类型
}

// 把ExpType转换为sql表达式，把LogicType转换为sql使用字符
func (column *Column) convert() error {
	if column.ExpType == "" {
		column.ExpType = Eq
	}
	if v, ok := ExpMap[strings.ToLower(column.ExpType)]; ok {
		column.ExpType = v
		if column.ExpType == " LIKE " {
			column.Value = fmt.Sprintf("%%%v%%", column.Value)
		}
	} else {
		return fmt.Errorf("unknown column expression type '%s'", column.ExpType)
	}

	if column.LogicType == "" {
		column.LogicType = AND
	}
	if v, ok := logicMap[strings.ToLower(column.LogicType)]; ok {
		column.LogicType = v
	} else {
		return fmt.Errorf("unknown logic type '%s'", column.LogicType)
	}

	return nil
}

// GetQueryConditions 获取查询条件，如果只查询一列，忽略第一列的逻辑类型
func GetQueryConditions(columns []*Column) (string, []interface{}, error) {
	str := ""
	args := []interface{}{}
	if len(columns) == 0 {
		return str, nil, nil
	}

	isUseIN := true
	if len(columns) == 1 {
		isUseIN = false
	}
	field := columns[0].Name

	for i, column := range columns {
		err := column.convert()
		if err != nil {
			return "", nil, err
		}

		if i == 0 { // 忽略第一列的逻辑类型
			str = column.Name + column.ExpType + "?"
		} else {
			str += column.LogicType + column.Name + column.ExpType + "?"
		}
		args = append(args, column.Value)

		if isUseIN {
			if field != column.Name {
				isUseIN = false
				continue
			}
			if column.ExpType != ExpMap[Eq] {
				isUseIN = false
			}
		}
	}

	if isUseIN {
		str = field + " IN (?)"
		args = []interface{}{args}
	}

	return str, args, nil
}

func getExpsAndLogics(keyLen int, paramSrc string) ([]string, []string) {
	exps, logics := []string{}, []string{}
	param := strings.Replace(paramSrc, " ", "", -1)
	sps := strings.SplitN(param, "?", 2)
	if len(sps) == 2 {
		param = sps[1]
	}

	num := keyLen
	if num == 0 {
		return exps, logics
	}

	fields := []string{}
	kvs := strings.Split(param, "&")
	for _, kv := range kvs {
		if strings.Contains(kv, "page=") || strings.Contains(kv, "size=") || strings.Contains(kv, "sort=") {
			continue
		}
		fields = append(fields, kv)
	}

	// 根据不重复的key分为num组，在每组中判断exp和logic是否存在
	group := map[string]string{}
	for _, field := range fields {
		split := strings.SplitN(field, "=", 2)
		if len(split) != 2 {
			continue
		}

		if _, ok := group[split[0]]; ok {
			// 在一组中，如果exp不存在则填充默认值空，logic不存在则填充充默认值空
			exps = append(exps, group["exp"])
			logics = append(logics, group["logic"])

			group = map[string]string{}
			continue
		} else {
			group[split[0]] = split[1]
		}
	}

	// 处理最后一组
	exps = append(exps, group["exp"])
	logics = append(logics, group["logic"])

	return exps, logics
}

// GetColumns 通过参数获取列信息
func GetColumns(keys []string, values []interface{}, exps []string, logics []string, paramSrc string) ([]*Column, error) {
	keyLen := len(keys)
	if keyLen == 0 && len(values) == 0 {
		return nil, nil
	}
	if len(values) != keyLen {
		return nil, errors.New("len(values) must be equal to len(keys)")
	}
	if paramSrc == "" {
		if len(exps) != keyLen {
			return nil, errors.New("len(exps) must be equal to len(keys)")
		}
		if len(logics) != keyLen {
			return nil, errors.New("len(logics) must be equal to len(keys)")
		}
	} else {
		exps, logics = getExpsAndLogics(keyLen, paramSrc)
		if len(exps) != keyLen || len(logics) != keyLen {
			return nil, errors.New("param format is illegal")
		}
	}

	columns := make([]*Column, 0, keyLen)
	for i := 0; i < keyLen; i++ {
		columns = append(columns, &Column{
			Name:      keys[i],
			Value:     values[i],
			ExpType:   exps[i],
			LogicType: logics[i],
		})
	}

	return columns, nil
}
