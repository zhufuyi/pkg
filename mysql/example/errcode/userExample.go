package errcode

import "github.com/zhufuyi/pkg/gin/errcode"

const (
	// 需要修改字段
	userExampleName = "用户" // userExample对应的名称
	userExampleNO   = 1    // 每个资源名称对应唯一编号，值范围建议1~1000
)

// 服务级别错误码，有Err前缀
var (
	ErrCreateUserExample = errcode.NewError(genCode(userExampleNO)+1, "创建"+userExampleName+"失败") // 错误码20101
	ErrDeleteUserExample = errcode.NewError(genCode(userExampleNO)+2, "删除"+userExampleName+"失败")
	ErrUpdateUserExample = errcode.NewError(genCode(userExampleNO)+3, "更新"+userExampleName+"失败")
	ErrGetUserExample    = errcode.NewError(genCode(userExampleNO)+4, "获取"+userExampleName+"失败")
	// 每添加一个错误码，在上一个错误码基础上+1
)
