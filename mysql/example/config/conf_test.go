package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

func TestParseYAML(t *testing.T) {
	err := Init("conf.yml") // 解析yaml文件
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(ShowConfig())
}

// 测试更新配置文件
func TestWatch(t *testing.T) {
	err := Init("conf.yml")
	if err != nil {
		t.Error(err)
		return
	}

	fs := []func(){
		func() {
			fmt.Println("更新字段1")
		},
		func() {
			fmt.Println("更新字段2")
		},
	}
	WatchConfig(fs...)

	for i := 0; i < 30; i++ {
		fmt.Println("port:", ShowConfig().ServerPort)
		time.Sleep(time.Second)
	}
}
