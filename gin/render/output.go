package render

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/json-iterator/go"
	"github.com/zhufuyi/pkg/gin/errcode"
)

// JSONResponse 输出格式
type JSONResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
	//Data interface{} `json:"data,omitempty"`
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

func writeJSON(c *gin.Context, code int, res interface{}) {
	c.Writer.WriteHeader(code)
	writeContentType(c.Writer, jsonContentType)
	//err := json.NewEncoder(c.Writer).Encode(res)
	err := jsoniter.NewEncoder(c.Writer).Encode(res)
	if err != nil {
		fmt.Printf("json encode error, err = %s", err.Error())
	}
}

// JSON 输出JSONResponse，只取第一项返回
func JSON(c *gin.Context, code int, msg string, data ...interface{}) {
	resp := JSONResponse{Code: code}

	if msg != "" {
		resp.Msg = msg
	}

	if len(data) == 1 {
		resp.Data = data[0]
	}

	// 保证返回时data字段不为nil
	if resp.Data == nil {
		resp.Data = &struct{}{}
	}

	if c.IsAborted() {
		writeJSON(c, code, resp)
	} else {
		writeJSON(c, http.StatusOK, resp)
	}
	//writeJSON(c, code, resp)
}

// OK 正确200输出
func OK(c *gin.Context, data ...interface{}) {
	JSON(c, http.StatusOK, "ok", data...)
}

// Err 输出错误
func Err(c *gin.Context, code int, msg ...string) {
	if len(msg) > 0 {
		JSON(c, code, strings.Join(msg, ", "))
	} else {
		JSON(c, code, "")
	}
}

// Err400 无效参数
func Err400(c *gin.Context) {
	JSON(c, http.StatusBadRequest, "参数无效")
}

// Err400Msg 无效参数
func Err400Msg(c *gin.Context, msg interface{}) {
	JSON(c, http.StatusBadRequest, "参数无效, "+fmt.Sprint(msg))
}

// Err401 鉴权失败
func Err401(c *gin.Context) {
	JSON(c, http.StatusUnauthorized, "鉴权失败")
}

// Err403 禁止访问
func Err403(c *gin.Context) {
	JSON(c, http.StatusForbidden, "禁止访问")
}

// Err404 资源不存在
func Err404(c *gin.Context) {
	JSON(c, http.StatusNotFound, "资源不存在")
}

// Err408 请求超时
func Err408(c *gin.Context) {
	JSON(c, http.StatusRequestTimeout, "请求超时")
}

// Err409 资源冲突，已存在
func Err409(c *gin.Context) {
	JSON(c, http.StatusConflict, "资源冲突，已存在")
}

// Err410 资源消失
func Err410(c *gin.Context) {
	JSON(c, http.StatusGone, "资源已消失")
}

// Err500 服务内部错误
func Err500(c *gin.Context, msg interface{}) {
	JSON(c, http.StatusInternalServerError, fmt.Sprint(msg))
}

// Abort 中断并报错
func Abort(c *gin.Context, code int, msg string) {
	render.WriteJSON(c.Writer, JSONResponse{Code: code, Msg: msg})
	c.AbortWithStatus(http.StatusOK)
}

// RespondJSON 根据code返回对应结果
func RespondJSON(c *gin.Context, code int, msg ...interface{}) {
	switch code {
	case http.StatusOK:
		OK(c)
	case http.StatusBadRequest:
		Err400(c)
	case http.StatusUnauthorized:
		Err401(c)
	case http.StatusForbidden:
		Err403(c)
	case http.StatusNotFound:
		Err404(c)
	case http.StatusRequestTimeout:
		Err408(c)
	case http.StatusConflict:
		Err409(c)
	case http.StatusGone:
		Err410(c)
	case http.StatusInternalServerError:
		Err500(c, msg)
	}
}

// Success 正确
func Success(c *gin.Context, data ...interface{}) {
	JSON(c, 0, "ok", data...)
}

// Error 错误
func Error(c *gin.Context, err *errcode.Error, data ...interface{}) {
	JSON(c, err.Code(), err.Msg(), data...)
}
