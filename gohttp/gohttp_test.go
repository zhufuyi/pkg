package gohttp

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

type myBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

var requestAddr string

func init() {
	port, _ := getAvailablePort()
	requestAddr = fmt.Sprintf("http://localhost:%d", port)
	addr := fmt.Sprintf(":%d", port)

	r := gin.Default()
	oKFun := func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Printf("request parameters: uid=%s\n", uid)
		c.JSON(200, JSONResponse{
			Code: 0,
			Msg:  "ok",
			Data: fmt.Sprintf("uid=%v", uid),
		})
	}
	errFun := func(c *gin.Context) {
		uid := c.Query("uid")
		fmt.Printf("request parameters: uid=%s\n", uid)
		c.JSON(401, JSONResponse{
			Code: 401,
			Msg:  "authentication failure",
			Data: fmt.Sprintf("uid=%v", uid),
		})
	}

	oKPFun := func(c *gin.Context) {
		var body myBody
		c.BindJSON(&body)
		fmt.Println("body data:", body)
		c.JSON(200, JSONResponse{
			Code: 0,
			Msg:  "ok",
			Data: nil,
		})
	}
	errPFun := func(c *gin.Context) {
		var body myBody
		c.BindJSON(&body)
		fmt.Println("body data:", body)
		c.JSON(401, JSONResponse{
			Code: 401,
			Msg:  "authentication failure",
			Data: nil,
		})
	}

	r.GET("/getStandard", oKFun)
	r.GET("/getStandard_err", errFun)
	r.DELETE("/deleteStandard", oKFun)
	r.DELETE("/deleteStandard_err", errFun)
	r.POST("/postStandard", oKPFun)
	r.POST("/postStandard_err", errPFun)
	r.PUT("/putStandard", oKPFun)
	r.PUT("/putStandard_err", errPFun)
	r.PATCH("/patchStandard", oKPFun)
	r.PATCH("/patchStandard_err", errPFun)

	r.GET("/get", oKFun)
	r.GET("/get_err", errFun)
	r.DELETE("/delete", oKFun)
	r.DELETE("/delete_err", errFun)
	r.POST("/post", oKPFun)
	r.POST("/post_err", errPFun)
	r.PUT("/put", oKPFun)
	r.PUT("/put_err", errPFun)
	r.PATCH("/patch", oKPFun)
	r.PATCH("/patch_err", errPFun)

	go func() {
		err := r.Run(addr)
		if err != nil {
			panic(err)
		}
	}()
}

// 获取可用端口
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

// ------------------------------------------------------------------------------------------

func TestGetStandard(t *testing.T) {
	type args struct {
		url    string
		params KV
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONResponse
		wantErr bool
	}{
		{
			name: "get standard success",
			args: args{
				url:    requestAddr + "/getStandard",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
			wantErr: false,
		},
		{
			name: "get standard err",
			args: args{
				url:    requestAddr + "/getStandard_err",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "get not found",
			args: args{
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetStandard(tt.args.url, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteStandard(t *testing.T) {
	type args struct {
		url    string
		params KV
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONResponse
		wantErr bool
	}{
		{
			name: "delete standard success",
			args: args{
				url:    requestAddr + "/deleteStandard",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
			wantErr: false,
		},
		{
			name: "delete standard err",
			args: args{
				url:    requestAddr + "/deleteStandard_err",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "delete not found",
			args: args{
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeleteStandard(tt.args.url, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostStandard(t *testing.T) {
	type args struct {
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONResponse
		wantErr bool
	}{
		{
			name: "post standard success",
			args: args{
				url: requestAddr + "/postStandard",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "post standard error",
			args: args{
				url: requestAddr + "/postStandard_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				url: requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PostStandard(tt.args.url, tt.args.body, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutStandard(t *testing.T) {
	type args struct {
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONResponse
		wantErr bool
	}{
		{
			name: "put standard success",
			args: args{
				url: requestAddr + "/putStandard",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "put standard error",
			args: args{
				url: requestAddr + "/putStandard_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "put not found",
			args: args{
				url: requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PutStandard(tt.args.url, tt.args.body, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPatchStandard(t *testing.T) {
	type args struct {
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name    string
		args    args
		want    *JSONResponse
		wantErr bool
	}{
		{
			name: "patch standard success",
			args: args{
				url: requestAddr + "/patchStandard",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "patch standard error",
			args: args{
				url: requestAddr + "/patchStandard_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "patch not found",
			args: args{
				url: requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			want: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PatchStandard(tt.args.url, tt.args.body, tt.args.params...)
			if (err != nil) != tt.wantErr {
				t.Errorf("PostStandard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostStandard() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ------------------------------------------------------------------------------------------

func TestGet(t *testing.T) {
	type args struct {
		result interface{}
		url    string
		params KV
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult *JSONResponse
	}{
		{
			name: "get success",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/get",
				params: KV{"uid": 123},
			},
			wantErr: false,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
		},
		{
			name: "get err",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/get_err",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
		{
			name: "get not found",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Get(tt.args.result, tt.args.url, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*JSONResponse).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		result interface{}
		url    string
		params KV
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantResult *JSONResponse
	}{
		{
			name: "delete success",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/delete",
				params: KV{"uid": 123},
			},
			wantErr: false,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: "uid=123",
			},
		},
		{
			name: "delete err",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/delete_err",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
		{
			name: "delete not found",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/notfound",
				params: KV{"uid": 123},
			},
			wantErr: true,
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.result, tt.args.url, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*JSONResponse).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPost(t *testing.T) {
	type args struct {
		result interface{}
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name       string
		args       args
		wantResult *JSONResponse
		wantErr    bool
	}{
		{
			name: "post success",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/post",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "post error",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/post_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Post(tt.args.result, tt.args.url, tt.args.body, tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*JSONResponse).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPut(t *testing.T) {
	type args struct {
		result interface{}
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name       string
		args       args
		wantResult *JSONResponse
		wantErr    bool
	}{
		{
			name: "put success",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/put",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "put error",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/put_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Put(tt.args.result, tt.args.url, tt.args.body, tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*JSONResponse).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}

func TestPatch(t *testing.T) {
	type args struct {
		result interface{}
		url    string
		body   interface{}
		params []KV
	}
	tests := []struct {
		name       string
		args       args
		wantResult *JSONResponse
		wantErr    bool
	}{
		{
			name: "patch success",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/patch",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "ok",
				Data: nil,
			},
			wantErr: false,
		},
		{
			name: "patch error",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/patch_err",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
		{
			name: "post not found",
			args: args{
				result: &JSONResponse{},
				url:    requestAddr + "/notfound",
				body: &myBody{
					Name:  "foo",
					Email: "bar@gmail.com",
				},
			},
			wantResult: &JSONResponse{
				Code: 0,
				Msg:  "",
				Data: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Patch(tt.args.result, tt.args.url, tt.args.body, tt.args.params...); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.args.result.(*JSONResponse).Msg != tt.wantResult.Msg {
				t.Errorf("gotResult = %v, wantResult =  %v", tt.args.result, tt.wantResult)
			}
		})
	}
}
