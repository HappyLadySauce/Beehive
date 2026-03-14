package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveContactLogic {
	return &RemoveContactLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveContactLogic) RemoveContact(in *pb.RemoveContactRequest) (*pb.RemoveContactResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.RemoveContactResponse{}, nil
}
