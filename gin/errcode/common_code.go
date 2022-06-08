package errcode

var (
	Success                   = NewError(0, "成功")
	InvalidParams             = NewError(40000000, "入参错误")
	ServerError               = NewError(40000001, "服务内部错误")
	NotFound                  = NewError(40000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(40000003, "鉴权失败，找不到对应的AppKey和AppSecret")
	UnauthorizedTokenError    = NewError(40000004, "鉴权失败，Token错误")
	UnauthorizedTokenTimeout  = NewError(40000005, "鉴权失败，Token超时")
	UnauthorizedTokenGenerate = NewError(40000006, "鉴权失败，Token生成失败")
	TooManyRequests           = NewError(40000007, "请求过多")
	AlreadyExists             = NewError(40000008, "已存在")
)
