package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-friend/friend"
	"github.com/HappyLadySauce/Beehive/app/beehive-friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友列表
func (l *GetFriendsLogic) GetFriends(in *friend.GetFriendsRequest) (*friend.FriendsResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.FriendsResponse{}, nil
}
