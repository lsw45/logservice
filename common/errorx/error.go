package errorx

import (
	"fmt"
)

type CodeError struct {
	errCode ErrInt
	errMsg  string
}

//返回给前端的错误码
func (e *CodeError) GetErrCode() ErrInt {
	return e.errCode
}

//返回给前端显示端错误信息
func (e *CodeError) GetErrMsg() string {
	return e.errMsg
}

func (e *CodeError) Error() string {
	return fmt.Sprintf("ErrCode:%d，ErrMsg:%s", e.errCode, e.errMsg)
}

func NewErrCodeMsg(errCode ErrInt, errMsg string) *CodeError {
	return &CodeError{errCode: errCode, errMsg: errMsg}
}
func NewErrCode(errCode ErrInt) *CodeError {
	return &CodeError{errCode: errCode, errMsg: MapErrMsg(errCode)}
}

func NewErrMsg(errMsg string) *CodeError {
	return &CodeError{errCode: SERVER_COMMON_ERROR, errMsg: errMsg}
}
