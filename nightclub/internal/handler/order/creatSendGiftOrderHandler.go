package order

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/order"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func CreatSendGiftOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatSendGiftOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := order.NewCreatSendGiftOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreatSendGiftOrder(&req)
		result.HttpResult(r, w, resp, err)
	}
}
