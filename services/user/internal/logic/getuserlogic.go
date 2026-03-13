package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserResponse{}, nil
}
