package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/jwt"
)

var uid = "123"

func initServer2() {
	jwt.Init()

	addr := getAddr()
	r := gin.Default()

	tokenFun := func(c *gin.Context) {
		token, _ := jwt.GenerateToken(uid)
		fmt.Println("token =", token)
		c.String(200, token)
	}

	userFun := func(c *gin.Context) {
		c.JSON(200, "hello "+uid)
	}

	r.GET("/token", tokenFun)
	r.GET("/user/:id", JWT(), userFun) // 需要鉴权

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

func TestJWT(t *testing.T) {
	initServer2()

	token, err := get(requestAddr + "/token")
	if err != nil {
		t.Fatal(err)
	}

	authentication := fmt.Sprintf("Bearer %s", token)
	result, err := getUser(authentication)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}

func getUser(authentication string) (string, error) {
	client := &http.Client{}
	url := requestAddr + "/user/" + uid
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("Authentication", authentication)
	reqest.Header.Add("X-Uid", uid)
	if err != nil {
		return "", err
	}
	response, _ := client.Do(reqest)
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
