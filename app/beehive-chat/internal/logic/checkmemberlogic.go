package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-chat/chat"
	"github.com/HappyLadySauce/Beehive/app/beehive-chat/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckMemberLogic {
	return &CheckMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查用户是否在会话中
func (l *CheckMemberLogic) CheckMember(in *chat.CheckMemberRequest) (*chat.CheckMemberResponse, error) {
	// todo: add your logic here and delete this line

	return &chat.CheckMemberResponse{}, nil
}
