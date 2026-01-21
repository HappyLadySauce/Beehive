package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConversationInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConversationInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConversationInfoLogic {
	return &UpdateConversationInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新会话信息
func (l *UpdateConversationInfoLogic) UpdateConversationInfo(in *chat.UpdateConversationInfoRequest) (*chat.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.CommonResponse{}, nil
}
