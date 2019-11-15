package gohttp

import (
	"fmt"
	"testing"
)

func TestHTTPGetJson(t *testing.T) {
	url := "http://localhost:20080/user"
	params := map[string]interface{}{"name": "张三"}

	resp, err := HTTPGetJson(url, params)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(*resp)
}

func TestHTTPPostJson(t *testing.T) {
	url := "http://localhost:20080/login"
	body := struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}{"张三", "123456"}

	resp, err := HTTPPostJson(url, &body)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(*resp)
}
