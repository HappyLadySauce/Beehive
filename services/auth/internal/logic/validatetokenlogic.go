package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ValidateTokenLogic) ValidateToken(in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if in.GetAccessToken() == "" {
		return &pb.ValidateTokenResponse{
			Valid: false,
			UserId: "",
		}, nil
	}

	userID, _, _, err := loadAndTouchToken(l.ctx, l.svcCtx.Redis, in.GetAccessToken())
	if err != nil {
		return nil, err
	}
	if userID == "" {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
