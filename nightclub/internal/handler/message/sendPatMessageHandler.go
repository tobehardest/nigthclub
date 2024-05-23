package message

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/message"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func SendPatMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SendPatMessageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := message.NewSendPatMessageLogic(r.Context(), svcCtx)
		resp, err := l.SendPatMessage(&req)
		result.HttpResult(r, w, resp, err)
	}
}
