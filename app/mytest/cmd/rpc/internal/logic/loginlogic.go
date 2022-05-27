package logic

import (
	"context"

	"fmt"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"looklook/app/mytest/cmd/rpc/internal/svc"
	"looklook/app/mytest/cmd/rpc/mytest"
	"looklook/app/mytest/model"
	"looklook/common/tool"
	"looklook/common/xerr"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}
var ErrUserNoExistsError = xerr.NewErrMsg("用户不存在")
var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")
var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *mytest.LoginReq) (*mytest.LoginResp, error) {

	var userId int64
	var err error
	fmt.Println(in.AuthType)
	switch in.AuthType {
	case model.UserAuthTypeSystem:
		userId, err = l.loginByMobile(in.AuthKey, in.Password)
	default:
		return nil, xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}

	fmt.Println("aaaaaaaaaa",userId,err)
	if err != nil {
		return nil, err
	}
	fmt.Println("bbbbbbbb")
	//2、Generate the token, so that the service doesn't call rpc internally
	generateTokenLogic :=NewGenerateTokenLogic(l.ctx,l.svcCtx)
	tokenResp,err:=generateTokenLogic.GenerateToken(&mytest.GenerateTokenReq{
		UserId: userId,
	})
	fmt.Println("ccc",err)
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}

	fmt.Println("dddddddddddddd")
	return &mytest.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

func (l *LoginLogic) loginByMobile(mobile, password string) (int64, error) {

	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx,mobile)
	fmt.Println(8888888,user,err)
	if err != nil && err != model.ErrNotFound {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "根据手机号查询用户信息失败，mobile:%s,err:%v", mobile, err)
	}
	if user == nil {
		return 0, errors.Wrapf(ErrUserNoExistsError, "mobile:%s", mobile)
	}

	if !(tool.Md5ByString(password) == user.Password) {
		return 0, errors.Wrap(ErrUsernamePwdError, "密码匹配出错")
	}

	return user.Id, nil
}

func (l *LoginLogic) loginBySmallWx() error {
	return nil
}
