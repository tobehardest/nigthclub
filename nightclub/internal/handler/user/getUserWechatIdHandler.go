package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"nightclub/common/result"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetUserWechatIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserWechatIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetUserWechatIdLogic(r.Context(), svcCtx)
		resp, err := l.GetUserWechatId(&req)
		result.HttpResult(r, w, resp, err)
	}
}
