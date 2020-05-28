package gohttp

import (
	"fmt"
	"testing"

	"github.com/k0kubun/pp"
)

func TestHTTPGetJson(t *testing.T) {
	//url := "http://localhost:20080/user"
	url := "https://www.baidu.com"
	params := map[string]interface{}{"name": "张三"}

	result, err := HTTPGetJSON(url, params)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(*result)
}

func TestHTTPGet(t *testing.T) {
	type Result struct {
		Code int                 `json:"code"`
		Msg  string              `json:"msg"`
		Data map[string][]string `json:"data,omitempty"`
	}

	result := &Result{}
	url := "http://xxx.com/router"

	err := GetJSON(result, url, nil)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(result)
}

func TestHTTPPostJson(t *testing.T) {
	url := "http://localhost:20080/login"
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{"张三", "123456"}

	resp, err := HTTPPostJSON(url, &body)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(*resp)
}
