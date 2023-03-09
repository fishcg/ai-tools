package helper

const (
	CodeOK = 0

	// 1. 系统错误
	CodeUnknown = 100010001
	CodeBadCall = 100020001

	// 2. 参数错误
	CodeEmptyParam   = 200010001 // 参数为空
	CodeInvalidParam = 200010002 // 参数非法

	// 3. 数据错误
	CodeNotFound       = 300010001 // 数据不存在
	CodeDuplicateEntry = 300020001 // 数据重复插入或绑定
)
