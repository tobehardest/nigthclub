package chat

import (
	"net/http"
	"nightclub/common/result"

	"github.com/zeromicro/go-zero/rest/httpx"
	"nightclub/nightclub/internal/logic/chat"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetHasNewFriendChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetHasNewFriendChatReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := chat.NewGetHasNewFriendChatLogic(r.Context(), svcCtx)
		resp, err := l.GetHasNewFriendChat(&req)
		result.HttpResult(r, w, resp, err)
	}
}
