package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetConversationDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationDetailLogic {
	return &GetConversationDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话详情
func (l *GetConversationDetailLogic) GetConversationDetail(in *chat.GetConversationDetailRequest) (*chat.ConversationDetailResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.ConversationDetailResponse{}, nil
}
