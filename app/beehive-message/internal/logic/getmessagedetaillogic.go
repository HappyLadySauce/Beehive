package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-message/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-message/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMessageDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMessageDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMessageDetailLogic {
	return &GetMessageDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取消息详情
func (l *GetMessageDetailLogic) GetMessageDetail(in *message.GetMessageDetailRequest) (*message.MessageInfo, error) {
	// todo: add your logic here and delete this line

	return &message.MessageInfo{}, nil
}
