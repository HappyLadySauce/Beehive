package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-friend/friend"
	"github.com/HappyLadySauce/Beehive/app/beehive-friend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendRemarkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendRemarkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendRemarkLogic {
	return &UpdateFriendRemarkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新好友备注
func (l *UpdateFriendRemarkLogic) UpdateFriendRemark(in *friend.UpdateFriendRemarkRequest) (*friend.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &friend.CommonResponse{}, nil
}
