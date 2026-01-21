package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationsLogic {
	return &GetConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话列表
func (l *GetConversationsLogic) GetConversations(in *chat.GetConversationsRequest) (*chat.ConversationsResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.ConversationsResponse{}, nil
}
