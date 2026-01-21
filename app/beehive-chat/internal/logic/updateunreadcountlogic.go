package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUnreadCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUnreadCountLogic {
	return &UpdateUnreadCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新未读数
func (l *UpdateUnreadCountLogic) UpdateUnreadCount(in *chat.UpdateUnreadCountRequest) (*chat.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.CommonResponse{}, nil
}
