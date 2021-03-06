// Code generated by goctl. DO NOT EDIT!
// Source: users.proto

package server

import (
	"context"

	"looklook/app/mytest/cmd/rpc/internal/logic"
	"looklook/app/mytest/cmd/rpc/internal/svc"
	"looklook/app/mytest/cmd/rpc/pb/types/pb"
)

type MytestServer struct {
	svcCtx *svc.ServiceContext
	pb.UnimplementedMytestServer
}

func NewMytestServer(svcCtx *svc.ServiceContext) *MytestServer {
	return &MytestServer{
		svcCtx: svcCtx,
	}
}

func (s *MytestServer) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginResp, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

func (s *MytestServer) Register(ctx context.Context, in *pb.RegisterReq) (*pb.RegisterResp, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

func (s *MytestServer) GetUserInfo(ctx context.Context, in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

func (s *MytestServer) GetUserAuthByAuthKey(ctx context.Context, in *pb.GetUserAuthByAuthKeyReq) (*pb.GetUserAuthByAuthKeyResp, error) {
	l := logic.NewGetUserAuthByAuthKeyLogic(ctx, s.svcCtx)
	return l.GetUserAuthByAuthKey(in)
}

func (s *MytestServer) GetUserAuthByUserId(ctx context.Context, in *pb.GetUserAuthByUserIdReq) (*pb.GetUserAuthyUserIdResp, error) {
	l := logic.NewGetUserAuthByUserIdLogic(ctx, s.svcCtx)
	return l.GetUserAuthByUserId(in)
}

func (s *MytestServer) GenerateToken(ctx context.Context, in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	l := logic.NewGenerateTokenLogic(ctx, s.svcCtx)
	return l.GenerateToken(in)
}
