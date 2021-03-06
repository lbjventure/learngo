// Code generated by goctl. DO NOT EDIT!
// Source: users.proto

package mytest

import (
	"context"

	"looklook/app/mytest/cmd/rpc/pb/types/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GenerateTokenReq         = pb.GenerateTokenReq
	GenerateTokenResp        = pb.GenerateTokenResp
	GetUserAuthByAuthKeyReq  = pb.GetUserAuthByAuthKeyReq
	GetUserAuthByAuthKeyResp = pb.GetUserAuthByAuthKeyResp
	GetUserAuthByUserIdReq   = pb.GetUserAuthByUserIdReq
	GetUserAuthyUserIdResp   = pb.GetUserAuthyUserIdResp
	GetUserInfoReq           = pb.GetUserInfoReq
	GetUserInfoResp          = pb.GetUserInfoResp
	LoginReq                 = pb.LoginReq
	LoginResp                = pb.LoginResp
	RegisterReq              = pb.RegisterReq
	RegisterResp             = pb.RegisterResp
	User                     = pb.User
	UserAuth                 = pb.UserAuth

	Mytest interface {
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
		GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
		GetUserAuthByAuthKey(ctx context.Context, in *GetUserAuthByAuthKeyReq, opts ...grpc.CallOption) (*GetUserAuthByAuthKeyResp, error)
		GetUserAuthByUserId(ctx context.Context, in *GetUserAuthByUserIdReq, opts ...grpc.CallOption) (*GetUserAuthyUserIdResp, error)
		GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error)
	}

	defaultMytest struct {
		cli zrpc.Client
	}
)

func NewMytest(cli zrpc.Client) Mytest {
	return &defaultMytest{
		cli: cli,
	}
}

func (m *defaultMytest) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultMytest) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultMytest) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

func (m *defaultMytest) GetUserAuthByAuthKey(ctx context.Context, in *GetUserAuthByAuthKeyReq, opts ...grpc.CallOption) (*GetUserAuthByAuthKeyResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.GetUserAuthByAuthKey(ctx, in, opts...)
}

func (m *defaultMytest) GetUserAuthByUserId(ctx context.Context, in *GetUserAuthByUserIdReq, opts ...grpc.CallOption) (*GetUserAuthyUserIdResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.GetUserAuthByUserId(ctx, in, opts...)
}

func (m *defaultMytest) GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error) {
	client := pb.NewMytestClient(m.cli.Conn())
	return client.GenerateToken(ctx, in, opts...)
}
