package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RemoveContactLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveContactLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveContactLogic {
	return &RemoveContactLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *RemoveContactLogic) RemoveContact(in *pb.RemoveContactRequest) (*pb.RemoveContactResponse, error) {
	if in.GetOwnerId() == "" || in.GetContactUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner_id and contact_user_id are required")
	}
	if err := l.svcCtx.ContactMod.Remove(in.GetOwnerId(), in.GetContactUserId()); err != nil {
		return nil, status.Errorf(codes.Internal, "remove contact failed: %v", err)
	}
	return &pb.RemoveContactResponse{}, nil
}
