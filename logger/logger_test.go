package logger

import (
	"fmt"
	"testing"
)

type people struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func init() {
	//InitLogger(false, "", "debug") // 以console数据格式输出到控台(默认)
	//InitLogger(false, "", "debug", "json") // 以json数据格式输出到控台
	//InitLogger(true, "out.log", "debug") // 以json数据格式输出到文件
}

func TestLogger(t *testing.T) {
	Debug("this is debug")
	Info("this is info")
	Warn("this is warn")
	Error("this is error")

	p := &people{"张三", 11}
	ps := []people{{"张三", 11}, {"李四", 12}}
	pMap := map[string]people{"123": *p, "456": *p}
	Debug("this is debug object", Any("object1", p), Any("object2", ps), Any("object3", pMap))
	Info("err is not equal nil ", Any("object", ps))
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", String("string", "hello golang"))
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark type int", Int("int", i))
	}
}

func BenchmarkAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark type any", Any(fmt.Sprintf("object_%d", i), &people{"张三", 11}))
	}
}
