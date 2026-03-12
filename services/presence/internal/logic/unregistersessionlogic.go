package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnregisterSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnregisterSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnregisterSessionLogic {
	return &UnregisterSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnregisterSessionLogic) UnregisterSession(in *pb_presencepb.UnregisterSessionRequest) (*pb_presencepb.UnregisterSessionResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_presencepb.UnregisterSessionResponse{}, nil
}
