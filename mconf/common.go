package mconf

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// YamlFormat yaml格式
	YamlFormat = "yaml"
	// YamlFormat json格式
	JsonFormat = "json"
	// YamlFormat toml格式
	TomlFormat = "toml"

	// NotFound 记录yaml未找到的key的统一值
	NotFound = "__not_found__"
)

// Bytes2Str bytes转string
func Bytes2Str(val []byte) string {
	str := string(val)
	str = strings.Replace(str, "\n", "", -1)
	return strings.Trim(str, " ")
}

// Bytes2Int bytes转int
func Bytes2Int(val []byte) int {
	str := Bytes2Str(val)
	n, _ := strconv.Atoi(str)
	return n
}

// Bytes2Map bytes转map
func Bytes2Map(val []byte) map[string]string {
	str := string(val)
	str = strings.Trim(str, " ")
	ss := strings.Split(str, "\n")
	out := map[string]string{}

	for _, s := range ss {
		subStrs := strings.Split(s, ":")
		if len(subStrs) == 2 {
			out[strings.Trim(subStrs[0], " ")] = strings.Trim(subStrs[1], " ")
		}
	}

	return out
}

// Bytes2Slice bytes转slice
func Bytes2Slice(val []byte) []string {
	str := string(val)
	str = strings.Trim(str, " ")
	ss := strings.Split(str, "\n")
	out := []string{}
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// k8s的容器启动参数
// ---------------------------------------------------------------------------------------

// Bytes2SliceForArgs bytes转slice
func Bytes2SliceForArgs(val []byte) []string {
	str := Bytes2Str(val)

	strs := []string{}
	if strings.Contains(str, " -") {
		ss := strings.Split(str, " -")
		for i, s := range ss {
			s = strings.Trim(s, " ")
			if i == 0 {
				strs = append(strs, s)
				continue
			}
			strs = append(strs, "-"+s)
		}
	} else {
		return []string{str}
	}

	return strs
}

// ResourcesMap2Str 资源字典转字符串
func ResourcesMap2Str(m map[string]string) string {
	_, values := resourcesSizeMap2Slice(m)
	return strings.Join(values, ", ")
}

// 把map转为适合规则的slice
func resourcesSizeMap2Slice(kvs map[string]string) ([]string, []string) {
	valuesType := []string{}
	values := []string{}
	if v, ok := kvs["cpu"]; ok {
		valuesType = append(valuesType, "string")
		values = append(values, fmt.Sprintf("cpu=%s", v))
	}
	if v, ok := kvs["memory"]; ok {
		valuesType = append(valuesType, "string")
		values = append(values, fmt.Sprintf("memory=%s", v))
	}
	return valuesType, values
}

// 参数字典转字符串
func argsMap2Str(m map[string]string, sep string) string {
	str := ""
	for k, v := range m {
		k = "-" + trimKey(k)
		if v == "" {
			str = fmt.Sprintf("%s %s ", str, k)
		} else {
			if sep == "" {
				sep = " "
			}
			str = fmt.Sprintf("%s %s%s%s ", str, k, sep, v)
		}
	}

	return strings.Trim(str, " ")
}

// 把二维args转为一维度slice
func args2sliceAndType(args [][]string) []string {
	var out []string
	for _, arg := range args {
		out = append(out, strings.Join(arg, " "))
	}

	return out
}

// 字符串转map，a=1,b=2 --> {a:1,b:2}
func str2Map(str string) map[string]string {
	str = strings.Trim(str, " ")
	ss := strings.Split(str, ",")
	out := map[string]string{}

	for _, s := range ss {
		if strings.Contains(s, "=") {
			n := strings.Index(s, "=")
			key := strings.Trim(s[:n], " ")
			value := strings.Trim(s[n+1:], " ")
			fmt.Println(key, value)
			out[key] = value
		} else {
			out[strings.Trim(s, " ")] = ""
		}

	}

	return out
}

// ---------------------------------------------------------------------------------------

// PutRecord 已成功替换的参数记录
type PutRecord struct {
	Key    string `json:"key"`
	OldVal string `json:"oldVal"` // 如果存在则替换，如果查找的key不存在，则补充统一值__not_found__，表示是新添加
	NewVal string `json:"newVal"`
}

// 添加或替换参数
func addOrReplaceArgs(args [][]string, argsMap map[string]string) []*PutRecord {
	records := []*PutRecord{}

	addKVs := []string{}
	for i, arg := range args {
		for j, s := range arg {
			record, newKV := getMatchArgs(s, argsMap)
			if newKV == "" {
				continue
			}
			args[i][j] = newKV
			records = append(records, record)
		}
	}

	// 判断是否有需要新添加的kv
	addRecords := []*PutRecord{}
	if len(argsMap) > len(records) {
		for k, v := range argsMap {
			if !isExistKey(records, k) {
				sep := "="   // 默认分割符是=号
				if v == "" { // 如果值为空改为空格
					sep = ""
				}
				addKVs = append(addKVs, fmt.Sprintf(`-%s%s%s`, trimKey(k), sep, v))
				addRecords = append(addRecords, &PutRecord{
					Key:    k,
					OldVal: NotFound,
					NewVal: v,
				})
			}
		}
	}

	if len(addKVs) > 0 {
		args[len(args)-1] = append(args[len(args)-1], addKVs...) // 添加新kv
		records = append(records, addRecords...)                 // 添加新记录
	}

	return records
}

// 获取匹配的参数，如果匹配，返回替换后的kv值
func getMatchArgs(s string, argsMap map[string]string) (*PutRecord, string) {
	for k, v := range argsMap {
		kTmp := trimKey(k)
		sep := "" // 原参数k和v之间的分割符，只支持=号和空格两种

		if strings.Contains(s, kTmp) { // 判读是否匹配key(去掉前缀-号)
			index := strings.Index(s, kTmp)
			sepIndex := 0
			isSpace := false
			// 搜索定位分割符位置和值
			for n, char := range s[index+len(kTmp):] {
				if char == '=' {
					sep = "="
					sepIndex = index + len(kTmp) + n
					break
				}
				if char == ' ' {
					isSpace = true
					continue
				}
				if (char >= 'A' && char <= 'z') || (char >= '0' && char <= '9') || char == '_' {
					if isSpace {
						sep = " "
						sepIndex = index + len(kTmp) + n - 1
						break
					}
				}
			}

			if sep == "" {
				continue
			}

			// 验证分割后的key是否完全一致
			if kTmp != trimKey(s[0:sepIndex]) {
				continue
			}

			newKV := s[0:sepIndex] + sep + v // 保持原样
			replaceRecord := &PutRecord{
				Key:    k,
				OldVal: strings.Replace(s, s[0:sepIndex]+sep, "", -1),
				NewVal: v,
			}
			return replaceRecord, newKV
		}
	}

	return nil, ""
}

func trimKey(k string) string {
	kTmp := strings.Trim(k, " ")
	for len(kTmp) > 0 {
		if kTmp[0] != '-' {
			break
		}
		kTmp = strings.TrimLeft(kTmp, "-")
	}
	return kTmp
}

func isExistKey(records []*PutRecord, k string) bool {
	for _, record := range records {
		if record.Key == k {
			return true
		}
	}
	return false
}
