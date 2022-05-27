package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"looklook/app/mytest/cmd/api/internal/logic/user"
	"looklook/app/mytest/cmd/api/internal/svc"
	"looklook/app/mytest/cmd/api/internal/types"
	"looklook/common/result"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(req)
		result.HttpResult(r, w, resp, err)
	}
}
