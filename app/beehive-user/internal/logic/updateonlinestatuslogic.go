package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-user/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOnlineStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOnlineStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOnlineStatusLogic {
	return &UpdateOnlineStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新在线状态
func (l *UpdateOnlineStatusLogic) UpdateOnlineStatus(in *user.UpdateOnlineStatusRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
