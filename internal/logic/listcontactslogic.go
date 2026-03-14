package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListContactsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListContactsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactsLogic {
	return &ListContactsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListContactsLogic) ListContacts(in *pb.ListContactsRequest) (*pb.ListContactsResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListContactsResponse{}, nil
}
