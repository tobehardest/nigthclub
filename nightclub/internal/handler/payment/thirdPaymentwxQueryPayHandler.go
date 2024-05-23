package payment

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/payment"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func ThirdPaymentwxQueryPayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentWxPayQueryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewThirdPaymentwxQueryPayLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentwxQueryPay(&req)
		result.HttpResult(r, w, resp, err)
	}
}
