package result

import (
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
	"nightclub/common/xerr"
)

// http返回
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {

	if err == nil {
		// 成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		code := xerr.SERVER_COMMON_ERROR
		msg := "服务器开小差啦，稍后再来试一试"

		causeErr := errors.Cause(err)                // err类型
		if e, ok := causeErr.(*xerr.CodeError); ok { //自定义错误类型
			//自定义CodeError
			code = e.GetCode()
			msg = e.GetMsg()
		} else {
			if gstatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := uint32(gstatus.Code())
				if xerr.IsCodeErr(grpcCode) { //区分自定义错误跟系统底层、db等错误，底层、db错误不能返回给前端
					code = grpcCode
					msg = gstatus.Message()
				}
			}
		}

		logx.WithContext(r.Context()).Errorf("【API-ERR】 : %+v ", err)

		httpx.WriteJson(w, http.StatusBadRequest, Error(code, msg))
	}
}
