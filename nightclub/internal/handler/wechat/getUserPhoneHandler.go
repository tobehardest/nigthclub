package wechat

import (
	"net/http"
	"nightclub/common/result"
	"nightclub/nightclub/internal/logic/wechat"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserPhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPhoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := wechat.NewGetUserPhoneLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPhone(&req)
		result.HttpResult(r, w, resp, err)
	}
}
