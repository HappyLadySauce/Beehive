package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-message/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-message/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessagesLogic {
	return &GetMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取历史消息
func (l *GetMessagesLogic) GetMessages(in *message.GetMessagesRequest) (*message.MessagesResponse, error) {
	// todo: add your logic here and delete this line

	return &message.MessagesResponse{}, nil
}
