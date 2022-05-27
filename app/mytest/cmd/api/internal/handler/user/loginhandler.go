package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"looklook/app/mytest/cmd/api/internal/logic/user"
	"looklook/app/mytest/cmd/api/internal/svc"
	"looklook/app/mytest/cmd/api/internal/types"
	"looklook/common/result"

	//"github.com/zeromicro/go-zero/rest/httpx"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(req)
		result.HttpResult(r, w, resp, err)
	}
}
