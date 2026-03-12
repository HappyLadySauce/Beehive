package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshSessionLogic {
	return &RefreshSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshSessionLogic) RefreshSession(in *pb_presencepb.RefreshSessionRequest) (*pb_presencepb.RefreshSessionResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_presencepb.RefreshSessionResponse{}, nil
}
