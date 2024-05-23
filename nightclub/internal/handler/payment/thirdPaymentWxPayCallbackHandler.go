package payment

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"nightclub/nightclub/internal/logic/payment"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func ThirdPaymentWxPayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayCallbackReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := payment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentWxPayCallback(&req)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("【API-ERR】 ThirdPaymentWxPayCallbackHandler : %+v ", err)
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
		logx.Infof("Return : %v ", resp)
		fmt.Fprint(w.(http.ResponseWriter), resp)
	}
}
