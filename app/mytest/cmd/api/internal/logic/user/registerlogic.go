package user

import (
	"context"

	"looklook/app/mytest/cmd/api/internal/svc"
	"looklook/app/mytest/cmd/rpc/mytest"
	"looklook/app/mytest/model"
	"looklook/app/mytest/cmd/api/internal/types"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/pkg/errors"
)


type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) RegisterLogic {
	return RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req types.RegisterReq) (resp *types.RegisterResp, err error) {
	// todo: add your logic here and delete this line
	registerRsp,err:=l.svcCtx.MytestRpc.Register(l.ctx,&mytest.RegisterReq{
		AuthType: model.UserAuthTypeSystem,
		Mobile:   req.Mobile,
		AuthKey:  req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %+v", req)
	}
	_ =copier.Copy(&resp,registerRsp)

	return resp,nil
}
