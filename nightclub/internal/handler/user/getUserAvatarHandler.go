package user

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"
)

func GetUserAvatarHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserAvatarReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewGetUserAvatarLogic(r.Context(), svcCtx, w)
		err := l.GetUserAvatar(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
