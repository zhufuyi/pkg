package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

var requestAddr string

func init() {
	addr := getAddr()
	r := gin.Default()

	// 默认打印日志
	//	r.Use(Logging())

	// 自定义打印日志
	r.Use(Logging(
		WithMaxLen(400),
		WithRequestID(),
		//WithIgnoreRoutes("/hello"), // 忽略/hello
	))

	// 自定义zap log
	//log, _ := logger.Init(logger.WithFormat("json"))
	//r.Use(Logging(
	//	WithLog(log),
	//))

	pingFun := func(c *gin.Context) {
		c.JSON(200, "pong")
	}
	pongFun := func(c *gin.Context) {
		c.JSON(200, "ping")
	}
	helloFun := func(c *gin.Context) {
		c.JSON(200, "hello world")
	}

	r.GET("/ping", pingFun)
	r.GET("/pong", pongFun)

	r.GET("/hello", helloFun)
	r.DELETE("/hello", helloFun)
	r.POST("/hello", helloFun)
	r.PUT("/hello", helloFun)
	r.PATCH("/hello", helloFun)

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

// ------------------------------------------------------------------------------------------

func TestRequest(t *testing.T) {
	wantPong := `"pong"`
	wantPing := `"ping"`
	wantHello := `"hello world"`
	type User struct {
		Name string `json:"name"`
	}

	t.Run("ping", func(t *testing.T) {
		got, err := get(requestAddr + "/pong")
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantPing {
			t.Errorf("got: %s, want: %s", got, wantPing)
		}
	})

	t.Run("pong", func(t *testing.T) {
		got, err := get(requestAddr + "/ping")
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantPong {
			t.Errorf("got: %s, want: %s", got, wantPong)
		}
	})

	t.Run("get hello", func(t *testing.T) {
		got, err := get(requestAddr + "/hello")
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("delete hello", func(t *testing.T) {
		got, err := delete(requestAddr + "/hello")
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("post hello", func(t *testing.T) {
		got, err := post(requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("put hello", func(t *testing.T) {
		got, err := put(requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	t.Run("patch hello", func(t *testing.T) {
		got, err := patch(requestAddr+"/hello", &User{"foo"})
		if err != nil {
			t.Error(err)
			return
		}
		if string(got) != wantHello {
			t.Errorf("got: %s, want: %s", got, wantHello)
		}
	})

	time.Sleep(time.Millisecond * 100)
}

func getAddr() string {
	port, _ := getAvailablePort()
	requestAddr = fmt.Sprintf("http://localhost:%d", port)
	return fmt.Sprintf(":%d", port)
}

func getAvailablePort() (int, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:0", "0.0.0.0"))
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		return 0, err
	}

	port := listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()

	return port, err
}

func do(method string, url string, body interface{}) ([]byte, error) {
	var (
		resp        *http.Response
		err         error
		contentType = "application/json"
	)

	v, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodGet:
		resp, err = http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	case http.MethodPost:
		resp, err = http.Post(url, contentType, bytes.NewReader(v))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	case http.MethodDelete, http.MethodPut, http.MethodPatch:
		req, err := http.NewRequest(method, url, bytes.NewReader(v))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

	default:
		fmt.Errorf("%s method not supported", method)
	}

	return ioutil.ReadAll(resp.Body)
}

func get(url string) ([]byte, error) {
	return do(http.MethodGet, url, nil)
}

func delete(url string) ([]byte, error) {
	return do(http.MethodDelete, url, nil)
}

func post(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPost, url, body)
}

func put(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPut, url, body)
}

func patch(url string, body interface{}) ([]byte, error) {
	return do(http.MethodPatch, url, body)
}
