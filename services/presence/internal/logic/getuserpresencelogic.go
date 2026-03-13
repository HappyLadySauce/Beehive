package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPresenceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPresenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPresenceLogic {
	return &GetUserPresenceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPresenceLogic) GetUserPresence(in *pb.GetUserPresenceRequest) (*pb.GetUserPresenceResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserPresenceResponse{}, nil
}
