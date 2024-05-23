package xerr

import "fmt"

/**
常用通用固定错误
*/

type CodeError struct {
	Code uint32 `json:"code"`
	Msg  string `json:"msg"`
}

func (e *CodeError) GetCode() uint32 {
	return e.Code
}

func (e *CodeError) GetMsg() string {
	return e.Msg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.Code, e.Msg)
}

func NewErrCodeMsg(code uint32, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}
func NewErrCode(code uint32) *CodeError {
	return &CodeError{Code: code, Msg: MapErrMsg(code)}
}

func NewErrMsg(msg string) *CodeError {
	return &CodeError{Code: SERVER_COMMON_ERROR, Msg: msg}
}
