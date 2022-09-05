package errorx

import (
	"context"
	"log-ext/common/response"
	"strings"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

//api错误处理
func HandleErrorfunc(err error) (int, interface{}) {
	defaultCode := SERVER_COMMON_ERROR
	defaultMsg := "服务器开小差啦，稍后再来试一试"
	defaultRes := response.HttpResponse{Code: uint32(defaultCode), Msg: defaultMsg}
	causeErr := errors.Cause(err)           //api这边自定义的错误
	if e, ok := causeErr.(*CodeError); ok { //自定义错误类型
		logx.Errorf("【API-SRV-ERR】 %+v", err)
		return 412, response.HttpResponse{ //api返回的自定义CodeError
			Code: uint32(e.GetErrCode()),
			Msg:  e.GetErrMsg(),
			Data: struct{}{},
		}
	}
	if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
		grpcCode := ErrInt(gstatus.Code())
		if IsCodeErr(grpcCode) { //系统错误不返回，自定义错误透传
			if grpcCode == NOT_FOUNT_RETURN_NULL {
				return 200, response.HttpResponse{
					Code: 0,
					Msg:  "",
					Data: nil,
				}
			}
			return 412, response.HttpResponse{
				Code: uint32(grpcCode),
				Msg:  gstatus.Message(),
				Data: struct{}{},
			}
		}
		if grpcCode == 4 && strings.Contains(gstatus.Message(), "context deadline exceeded") {
			return 504, response.HttpResponse{
				Code: uint32(RPC_TIME_OUT),
				Msg:  message[RPC_TIME_OUT],
				Data: struct{}{},
			}
		}
		//非自定义grpc错误记录日志
		logx.Errorf("【API-SRV-ERR】 %+v", err)
		return 500, defaultRes
	}

	//处理go-zero预定义错误
	if errors.Is(err, context.DeadlineExceeded) { //超时
		logx.Errorf("【API-SRV-ERR】 %+v", err)
		return 504, response.HttpResponse{
			Code: uint32(DeadlineExceeded_ERROR),
			Msg:  MapErrMsg(DeadlineExceeded_ERROR),
			Data: struct{}{},
		}
	}

	//剩下的是api这边的系统错误，返回默认错误并打印日志，错误要wrap
	logx.Errorf("【API-SRV-ERR】 %+v", err)
	return 500, response.HttpResponse{
		Code: uint32(defaultCode),
		Msg:  defaultMsg,
		Data: struct{}{},
	}
}
