package user

import (
	"net/http"
	"nightclub/common/result"
	"nightclub/nightclub/internal/logic/user"
	"nightclub/nightclub/internal/svc"
	"nightclub/nightclub/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UserRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUserRegisterLogic(r, svcCtx)
		resp, err := l.UserRegister(&req)

		result.HttpResult(r, w, resp, err)
	}
}
