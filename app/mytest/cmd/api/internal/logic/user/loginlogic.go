package user

import (
	"context"

	"looklook/app/mytest/cmd/api/internal/svc"
	"looklook/app/mytest/cmd/api/internal/types"
	"looklook/app/mytest/cmd/rpc/mytest"
	"looklook/app/mytest/model"

	"github.com/jinzhu/copier"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) ( *types.LoginResp,  error) {
	// todo: add your logic here and delete this line

	loginResp, err := l.svcCtx.MytestRpc.Login(l.ctx, &mytest.LoginReq{
		AuthType: model.UserAuthTypeSystem,
		AuthKey:  req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var resp types.LoginResp
	_ = copier.Copy(&resp, loginResp)

	return &resp, nil
}
