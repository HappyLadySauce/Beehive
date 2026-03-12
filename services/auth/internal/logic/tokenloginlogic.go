package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type TokenLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTokenLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLoginLogic {
	return &TokenLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TokenLoginLogic) TokenLogin(in *pb_authpb.TokenLoginRequest) (*pb_authpb.LoginResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_authpb.LoginResponse{}, nil
}
