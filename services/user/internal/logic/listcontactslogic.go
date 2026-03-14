package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListContactsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListContactsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactsLogic {
	return &ListContactsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListContactsLogic) ListContacts(in *pb.ListContactsRequest) (*pb.ListContactsResponse, error) {
	if in.GetOwnerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_id is required")
	}
	ids, err := l.svcCtx.ContactMod.ListContactUserIDs(in.GetOwnerId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list contacts failed: %v", err)
	}
	return &pb.ListContactsResponse{ContactUserIds: ids}, nil
}
