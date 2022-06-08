package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type userForm struct {
	Name  string `json:"name" form:"name" binding:"required"`
	Age   int    `json:"age" form:"age" binding:"gte=0,lte=150"`
	Email string `json:"email" form:"email" binding:"email,required"`
}

func init() {
	r := gin.Default()
	binding.Validator = Gin()

	helloFun := func(c *gin.Context) {
		form := &userForm{}
		err := c.ShouldBindJSON(form)
		if err != nil {
			fmt.Println(err)
			c.JSON(400, "params is invalid")
			return
		}
		c.JSON(200, "hello world")
	}

	r.POST("/hello", helloFun)

	go func() {
		err := r.Run(":8080")
		if err != nil {
			panic(err)
		}
	}()
}

// -------------------------------------------------------------------------------------------

func post(url string, body interface{}) ([]byte, error) {
	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(v))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func TestValidate(t *testing.T) {
	wantHello := `"hello world"`
	wantParamErr := `"params is invalid"`

	t.Run("success", func(t *testing.T) {
		got, err := post("http://localhost:8080/hello", &userForm{
			Name:  "foo",
			Age:   10,
			Email: "bar@gmail.com",
		})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {

			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("missing field error", func(t *testing.T) {
		got, err := post("http://localhost:8080/hello", &userForm{
			Age:   10,
			Email: "bar@gmail.com",
		})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantParamErr {
			t.Errorf("got: %s, want: %s", got, wantParamErr)
		}
	})

	t.Run("field range  error", func(t *testing.T) {
		got, err := post("http://localhost:8080/hello", &userForm{
			Name:  "foo",
			Age:   -1,
			Email: "bar@gmail.com",
		})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantParamErr {
			t.Errorf("got: %s, want: %s", got, wantParamErr)
		}
	})

	t.Run("email error", func(t *testing.T) {
		got, err := post("http://localhost:8080/hello", &userForm{
			Name:  "foo",
			Age:   -10,
			Email: "bar",
		})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantParamErr {
			t.Errorf("got: %s, want: %s", got, wantParamErr)
		}
	})
}
