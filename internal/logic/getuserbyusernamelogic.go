package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByUsernameLogic) GetUserByUsername(in *pb.GetUserByUsernameRequest) (*pb.GetUserResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserResponse{}, nil
}
