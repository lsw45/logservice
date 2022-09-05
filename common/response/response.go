package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

//http响应
type HttpResponse struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OkJson(w http.ResponseWriter, v interface{}) {
	if v == nil {
		v = struct{}{}
	}
	httpx.OkJson(w, HttpResponse{Code: 0, Msg: "", Data: v})
}
