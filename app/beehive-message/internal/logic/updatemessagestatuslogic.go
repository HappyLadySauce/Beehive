package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-message/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-message/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMessageStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMessageStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMessageStatusLogic {
	return &UpdateMessageStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新消息状态
func (l *UpdateMessageStatusLogic) UpdateMessageStatus(in *message.UpdateMessageStatusRequest) (*message.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &message.CommonResponse{}, nil
}
