package user

import (
	"context"

	"looklook/app/mytest/cmd/api/internal/svc"
	"looklook/app/mytest/cmd/api/internal/types"
	"looklook/app/mytest/cmd/rpc/mytest"
	"looklook/common/ctxdata"

	"github.com/jinzhu/copier"


	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) DetailLogic {
	return DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req types.UserInfoReq) (resp *types.UserInfoResp, err error) {
	// todo: add your logic here and delete this line

	userId := ctxdata.GetUidFromCtx(l.ctx)
	userInfoResp,err := l.svcCtx.MytestRpc.GetUserInfo(l.ctx,&mytest.GetUserInfoReq{
    	Id:userId,
	})
	if err != nil {
		return nil, err
	}

	var userInfo types.User
	_ = copier.Copy(&userInfo, userInfoResp.User)

	return &types.UserInfoResp{
		UserInfo: userInfo,
	}, nil

}
