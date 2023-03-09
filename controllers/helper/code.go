package helper

const (
	CodeOK = 0

	// 1. 系统错误
	CodeUnknown = 100010001
	CodeBadCall = 100020001

	// 3. 参数错误
	CodeEmptyParam   = 200010001 // 参数为空
	CodeInvalidParam = 200010002 // 参数非法
)
