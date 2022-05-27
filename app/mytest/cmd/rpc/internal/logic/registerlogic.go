package logic

import (
	"context"

	"looklook/app/mytest/cmd/rpc/internal/svc"
	//"looklook/app/mytest/cmd/rpc/pb/types/pb"
	"looklook/common/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	//"fmt"
	"looklook/app/mytest/cmd/rpc/mytest"
	"looklook/app/mytest/model"
	"looklook/common/tool"
)
var ErrUserAlreadyRegisterError = xerr.NewErrMsg("user has been registered")
type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *mytest.RegisterReq) (*mytest.RegisterResp, error) {
	// todo: add your logic here and delete this line

	user,err:= l.svcCtx.UserModel.FindOneByMobile(l.ctx,in.Mobile)

	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "mobile:%s,err:%v", in.Mobile, err)
	}
	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v", in.Mobile, err)
	}

	var userId int64

	if err:=l.svcCtx.UserModel.Trans(l.ctx, func(ctx context.Context,session sqlx.Session) error {
		user :=new(model.Users)
		user.Mobile = in.Mobile
		if len(user.Nickname) == 0 {
			user.Nickname = tool.Krand( 8, tool.KC_RAND_KIND_ALL)
		}
		if len(in.Password) > 0 {
			user.Password = tool.Md5ByString(in.Password)
		}

		insertResult, err := l.svcCtx.UserModel.Insert(ctx,session, user)
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user Insert err:%v,user:%+v", err, user)
		}
		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user insertResult.LastInsertId err:%v,user:%+v", err, user)
		}
		userId = lastId

		userAuth := new(model.UserAuths)
		userAuth.UserId = lastId
		userAuth.AuthKey = in.AuthKey
		userAuth.AuthType = in.AuthType

		if _, err := l.svcCtx.UserAuthModel.Insert(ctx,session, userAuth); err != nil {
			return errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register db user_auth Insert err:%v,userAuth:%v", err, userAuth)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	//2、Generate the token, so that the service doesn't call rpc internally
	generateTokenLogic :=NewGenerateTokenLogic(l.ctx,l.svcCtx)
	tokenResp,err:=generateTokenLogic.GenerateToken(&mytest.GenerateTokenReq{
		UserId: userId,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}

	return &mytest.RegisterResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
