package logic

import (
	"context"
	"errors"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if in.GetAccessToken() == "" {
		return nil, errors.New("access_token is empty")
	}

	if l.svcCtx.Redis == nil {
		return nil, errors.New("redis client is nil")
	}

	if err := l.svcCtx.Redis.Del(l.ctx, tokenKey(in.GetAccessToken())).Err(); err != nil {
		return nil, err
	}

	return &pb.LogoutResponse{}, nil
}
