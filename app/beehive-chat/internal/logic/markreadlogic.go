package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkReadLogic {
	return &MarkReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 标记已读
func (l *MarkReadLogic) MarkRead(in *chat.MarkReadRequest) (*chat.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.CommonResponse{}, nil
}
