package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/app/beehive-user/internal/svc"
	"github.com/HappyLadySauce/Beehive/app/beehive-user/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCodeLogic {
	return &VerifyCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 验证邮箱验证码
func (l *VerifyCodeLogic) VerifyCode(in *user.VerifyCodeRequest) (*user.CommonResponse, error) {
	// todo: add your logic here and delete this line

	return &user.CommonResponse{}, nil
}
