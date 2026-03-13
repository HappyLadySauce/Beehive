package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserConversationsLogic {
	return &ListUserConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserConversationsLogic) ListUserConversations(in *pb.ListUserConversationsRequest) (*pb.ListUserConversationsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListUserConversationsResponse{}, nil
}
