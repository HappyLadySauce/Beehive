package logic

import (
	"context"
	"errors"

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

func (l *TokenLoginLogic) TokenLogin(in *pb.TokenLoginRequest) (*pb.LoginResponse, error) {
	if in.GetAccessToken() == "" {
		return nil, errors.New("access_token is empty")
	}

	userID, _, ttl, err := loadAndTouchToken(l.ctx, l.svcCtx.Redis, in.GetAccessToken())
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, errors.New("token invalid")
	}

	return &pb.LoginResponse{
		UserId:       userID,
		AccessToken:  in.GetAccessToken(),
		RefreshToken: "",
		ExpiresIn:    int64(ttl.Seconds()),
	}, nil
}
