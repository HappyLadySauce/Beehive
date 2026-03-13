package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMembersLogic {
	return &ListMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMembersLogic) ListMembers(in *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.ListMembersResponse{}, nil
}
