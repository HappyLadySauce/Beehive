package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-friend/friend"
	"github.com/HappyLadySauce/Beehive/app/beehive-friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendRequestsLogic {
	return &GetFriendRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取好友申请列表
func (l *GetFriendRequestsLogic) GetFriendRequests(in *friend.GetFriendRequestsRequest) (*friend.FriendRequestsResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.FriendRequestsResponse{}, nil
}
