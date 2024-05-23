package user

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func UpdateUserLocationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateUserLocationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateUserLocationLogic(r.Context(), svcCtx)
		resp, err := l.UpdateUserLocation(&req)
		result.HttpResult(r, w, resp, err)
	}
}
