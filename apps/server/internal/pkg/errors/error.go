package errors

import (
	stderrors "errors"
	"fmt"
)



type BusinessError struct {
	code int
	message string
}

// Error 实现 error 接口
func (e *BusinessError) Error() string {
	return fmt.Sprintf("[Code]:%d [Msg]:%s", e.code, e.message)
}


// GetCode 返回业务错误码
func (e *BusinessError) GetCode() int {
	return e.code
}


// GetMessage 返回业务错误消息
func (e *BusinessError) GetMessage() string {
	return e.message
}

// NewBusinessError 创建业务错误
func NewBusinessError(code int, message ...string) *BusinessError {
	msg := ""
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	} else {
		msg = GetErrorMessage(code)
	}

	return &BusinessError{
		code: code,
		message: msg,
	}
}

// Error 提供错误转换辅助方法
type ErrorHelper struct{}

// AsBusinessError 尝试把任意 error 转为 *BusinessError
func (e *ErrorHelper) AsBusinessError(err error) (*BusinessError, error) {
	var be *BusinessError
	if stderrors.As(err, &be) {
		return be, nil
	}
	return nil, err
}