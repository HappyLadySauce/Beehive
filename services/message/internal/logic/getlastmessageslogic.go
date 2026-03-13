package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/message/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/message/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLastMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLastMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLastMessagesLogic {
	return &GetLastMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLastMessagesLogic) GetLastMessages(in *pb.GetLastMessagesRequest) (*pb.GetLastMessagesResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetLastMessagesResponse{}, nil
}
