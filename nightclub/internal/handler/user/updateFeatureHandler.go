package user

import (
	"net/http"
	"nightclub/common/result"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UpdateFeatureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFeatureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateFeatureLogic(r.Context(), svcCtx)
		resp, err := l.UpdateFeature(&req)
		result.HttpResult(r, w, resp, err)
	}
}
