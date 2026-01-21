package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-friend/friend"
	"github.com/HappyLadySauce/Beehive/app/beehive-friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 检查是否为好友
func (l *IsFriendLogic) IsFriend(in *friend.IsFriendRequest) (*friend.IsFriendResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.IsFriendResponse{}, nil
}
