package user

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func UpdateWechatIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateWechatIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateWechatIdLogic(r.Context(), svcCtx)
		resp, err := l.UpdateWechatId(&req)
		result.HttpResult(r, w, resp, err)
	}
}
