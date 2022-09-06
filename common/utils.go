package common

import (
	"encoding/json"
	"log-ext/common/errorx"
	"log-ext/common/response"
)

func ParseBasicAuth(auth string) (key string, ok bool) {
	if auth == "" {
		return "", false
	}
	const prefix = "Bearer "
	if len(auth) < len(prefix) || !EqualFold(auth[:len(prefix)], prefix) {
		return "", false
	}
	return auth[len(prefix):], true
}

func EqualFold(s, t string) bool {
	if len(s) != len(t) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != t[i] {
			return false
		}
	}
	return true
}

func ThrowErr(err *errorx.CodeError) []byte {
	res := response.HttpResponse{
		Code: uint32(err.GetErrCode()),
		Msg:  err.GetErrMsg(),
		Data: struct{}{},
	}
	ErrBytes, _ := json.Marshal(res)
	return ErrBytes
}