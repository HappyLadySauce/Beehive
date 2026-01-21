package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-user/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
