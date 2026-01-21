package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-user/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyPasswordLogic {
	return &VerifyPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验密码
func (l *VerifyPasswordLogic) VerifyPassword(in *user.VerifyPasswordRequest) (*user.VerifyPasswordResponse, error) {
	// todo: add your logic here and delete this line

	return &user.VerifyPasswordResponse{}, nil
}
