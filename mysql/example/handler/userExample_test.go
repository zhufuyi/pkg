package handler

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jinzhu/copier"
	"github.com/zhufuyi/pkg/mysql"
	"github.com/zhufuyi/pkg/mysql/example/model"
)

func TestCopy1(t *testing.T) {
	form := &CreateUserExampleRequest{
		Name:     "lisi",
		Password: "123456",
		Email:    "foo@example.com",
		Phone:    "12345678901",
		Age:      11,
		Gender:   1,
	}

	userExample := &model.UserExample{}
	err := copier.Copy(userExample, form)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(userExample, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}

func TestCopy2(t *testing.T) {
	form := &UpdateUserExampleByIDRequest{
		ID: 817867015552241667,
		CreateUserExampleRequest: CreateUserExampleRequest{
			Name:     "lisi",
			Password: "123456",
			Email:    "foo@example.com",
			Phone:    "12345678901",
			Age:      11,
			Gender:   1,
		},
	}

	userExample := &model.UserExample{}
	err := copier.Copy(userExample, form)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(userExample, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}

func TestCopy3(t *testing.T) {
	userExample := &model.UserExample{
		Model: mysql.Model{
			ID:        817867015552241666,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:     "lisi",
		Password: "123456",
		Email:    "foo@example.com",
		Phone:    "12345678901",
		Age:      11,
		Gender:   1,
		Status:   2,
		LoginAt:  time.Now().Unix(),
	}
	data := &GetUserExampleByIDRespond{}
	err := copier.Copy(data, userExample)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}

func Test_convertUserExamples(t *testing.T) {
	userExamples := []model.UserExample{
		{
			Model: mysql.Model{
				ID:        817867015552241664,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:     "lisi",
			Password: "123456",
			Email:    "foo@example.com",
			Phone:    "12345678901",
			Age:      11,
			Gender:   1,
			Status:   2,
			LoginAt:  time.Now().Unix(),
		},
		{
			Model: mysql.Model{
				ID:        817867015552241665,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			Name:     "lisi1",
			Password: "1234561",
			Email:    "foo@example1.com",
			Phone:    "123456789011",
			Age:      12,
			Gender:   2,
			Status:   1,
			LoginAt:  time.Now().Unix(),
		},
	}

	data, err := convertUserExamples(userExamples)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(out))
}
