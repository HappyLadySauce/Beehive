package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateUserResponse{}, nil
}
