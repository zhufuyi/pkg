package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

// 系统级别错误码，无Err前缀
var (
	Success             = errcode.Success
	InvalidParams       = errcode.InvalidParams
	Unauthorized        = errcode.Unauthorized
	InternalServerError = errcode.InternalServerError
	NotFound            = errcode.NotFound
	AlreadyExists       = errcode.AlreadyExists
	Timeout             = errcode.Timeout
	TooManyRequests     = errcode.TooManyRequests
	Forbidden           = errcode.Forbidden
	LimitExceed         = errcode.LimitExceed
)

func genCode(NO int) int {
	return 20000 + NO*100
}
