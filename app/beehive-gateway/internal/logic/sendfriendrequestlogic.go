// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-gateway/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFriendRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 发送好友申请
func NewSendFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFriendRequestLogic {
	return &SendFriendRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendFriendRequestLogic) SendFriendRequest(req *types.SendFriendRequestReq) (resp *types.SendFriendRequestResp, err error) {
	// todo: add your logic here and delete this line

	return
}
