package helper

import (
	"net/http"
)

type HTTPError interface {
	Status() int
	Code() int64
	error
}

type ActionError struct {
	status  int
	code    int64
	message string
}

func (e *ActionError) Status() int {
	return e.status
}

func (e *ActionError) Code() int64 {
	return e.code
}

func (e *ActionError) Error() string {
	return e.message
}

func NewActionError(status int, code int64, msg string) *ActionError {
	return &ActionError{
		status,
		code,
		msg,
	}
}

func NewBadCallError(msg string) *ActionError {
	return &ActionError{
		http.StatusBadRequest,
		CodeBadCall,
		msg,
	}
}

// NewEmptyParamError 返回参数为空的错误
func NewEmptyParamError(msg string) *ActionError {
	return &ActionError{
		http.StatusBadRequest,
		CodeEmptyParam,
		msg,
	}
}

// NewInvalidParamError 返回参数不合法错误
func NewInvalidParamError(msg string) *ActionError {
	return &ActionError{
		http.StatusBadRequest,
		CodeInvalidParam,
		msg,
	}
}

// NewNotFoundError 返回元素不存在错误
func NewNotFoundError(msg string) *ActionError {
	return &ActionError{
		http.StatusNotFound,
		CodeNotFound,
		msg,
	}
}

func NewDuplicateEntry(msg string) *ActionError {
	return &ActionError{
		http.StatusBadRequest,
		CodeDuplicateEntry,
		msg,
	}
}

// 对报错进行规范
var (
	// 1. 系统错误
	ErrInternalServerError = NewActionError(http.StatusInternalServerError, CodeUnknown, "服务器忙，请稍后再试")

	// 2. 参数错误
	ErrInvalidParam = NewActionError(http.StatusBadRequest, CodeInvalidParam, "参数不合法")
)
