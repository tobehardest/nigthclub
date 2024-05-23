package common

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"nightclub/nightclub/internal/logic/common"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetSpecialEffectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetSpecialEffectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := common.NewGetSpecialEffectLogic(r.Context(), svcCtx, w)
		err := l.GetSpecialEffect(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
