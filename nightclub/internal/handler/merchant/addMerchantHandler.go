package merchant

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/merchant"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func AddMerchantHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddMerchantReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := merchant.NewAddMerchantLogic(r.Context(), svcCtx)
		resp, err := l.AddMerchant(&req)
		result.HttpResult(r, w, resp, err)
	}
}
