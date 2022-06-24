package model

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/zhufuyi/pkg/mysql"
)

var db *mysql.DB
var addr = "root:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func init() {
	var err error
	db, err = mysql.Init(addr, mysql.WithLog())
	if err != nil {
		panic(fmt.Sprintf("connect to mysql failed, err=%s, addr=%s", err, addr))
	}
}

func TestUser_Create(t *testing.T) {
	user := &UserExample{Name: "姜维", Age: 20, Gender: "男"}
	err := user.Create(db)
	if err != nil {
		t.Error(err)
	}

	if user.ID == 0 {
		t.Error("insert failed")
		return
	}

	t.Logf("id =%d", user.ID)
}

func TestUser_Updates(t *testing.T) {
	update := mysql.KV{"age": gorm.Expr("age  + ?", 1)}
	user := new(UserExample)
	err := user.Updates(db, update, "name = ?", "姜维")
	if err != nil {
		t.Error(err)
	}
}

func TestUser_Get(t *testing.T) {
	user := new(UserExample)
	err := user.Get(db, "name = ?", "姜维")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v", user)
}

func TestUser_Gets(t *testing.T) {
	user := new(UserExample)
	page := mysql.NewPage(0, 10, "-name")
	users, err := user.Gets(db, page, "")
	if err != nil {
		t.Error(err)
		return
	}
	for _, user := range users {
		t.Logf("%+v", user)
	}
}

func TestUser_Count(t *testing.T) {
	user := new(UserExample)
	count, err := user.Count(db, "id > ?", 10)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("count=%d", count)
}

func TestUser_TableName(t *testing.T) {
	var user UserExample
	fmt.Println(user.TableName())
}

func TestUser_Delete(t *testing.T) {
	user := new(UserExample)
	err := user.Delete(db, "name = ?", "姜维")
	if err != nil {
		t.Error(err)
	}
}
