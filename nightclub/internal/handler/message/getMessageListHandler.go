package message

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/message"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetMessageListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMessageListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := message.NewGetMessageListLogic(r.Context(), svcCtx)
		resp, err := l.GetMessageList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
