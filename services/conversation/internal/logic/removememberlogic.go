package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type RemoveMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveMemberLogic {
	return &RemoveMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveMemberLogic) RemoveMember(in *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	if in.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	_, err := l.svcCtx.Conv.FindByID(in.GetConversationId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "conversation not found")
		}
		l.Errorf("find conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "find conversation failed: %v", err)
	}
	if err := l.svcCtx.Conv.RemoveMember(in.GetConversationId(), in.GetUserId()); err != nil {
		l.Errorf("remove member failed: %v", err)
		return nil, status.Errorf(codes.Internal, "remove member failed: %v", err)
	}
	return &pb.RemoveMemberResponse{}, nil
}
