package payment

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"nightclub/nightclub/internal/logic/payment"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func ThirdPaymentwxPayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewThirdPaymentwxPayLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentwxPay(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
