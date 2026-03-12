package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOnlineSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineSessionsLogic {
	return &GetOnlineSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOnlineSessionsLogic) GetOnlineSessions(in *pb_presencepb.GetOnlineSessionsRequest) (*pb_presencepb.GetOnlineSessionsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_presencepb.GetOnlineSessionsResponse{}, nil
}
