package chat

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/chat"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetUserChatListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserChatListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewGetUserChatListLogic(r.Context(), svcCtx)
		resp, err := l.GetUserChatList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
