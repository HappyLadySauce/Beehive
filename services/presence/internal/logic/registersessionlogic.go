package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterSessionLogic {
	return &RegisterSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterSessionLogic) RegisterSession(in *pb_presencepb.RegisterSessionRequest) (*pb_presencepb.RegisterSessionResponse, error) {
	// todo: add your logic here and delete this line

	return &pb_presencepb.RegisterSessionResponse{}, nil
}
