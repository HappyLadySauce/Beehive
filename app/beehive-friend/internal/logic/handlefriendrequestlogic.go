package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-friend/friend"
	"github.com/HappyLadySauce/Beehive/app/beehive-friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type HandleFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleFriendRequestLogic {
	return &HandleFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 处理好友申请
func (l *HandleFriendRequestLogic) HandleFriendRequest(in *friend.HandleFriendRequestRequest) (*friend.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.CommonResponse{}, nil
}
