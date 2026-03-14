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

type DeclineJoinRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeclineJoinRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeclineJoinRequestLogic {
	return &DeclineJoinRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeclineJoinRequestLogic) DeclineJoinRequest(in *pb.DeclineJoinRequestRequest) (*pb.DeclineJoinRequestResponse, error) {
	if in.GetConversationId() == "" || in.GetRequestId() == "" || in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id, request_id and user_id are required")
	}
	member, err := l.svcCtx.Conv.GetMember(in.GetConversationId(), in.GetUserId())
	if err != nil || member == nil || member.Status != "active" {
		return nil, status.Error(codes.PermissionDenied, "not a member or not active")
	}
	if member.Role != "owner" && member.Role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "only owner or admin can decline")
	}
	if err := l.svcCtx.JoinReq.Decline(in.GetRequestId(), in.GetConversationId(), in.GetUserId()); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "request not found or not pending")
		}
		l.Errorf("decline join request failed: %v", err)
		return nil, status.Errorf(codes.Internal, "decline failed: %v", err)
	}
	return &pb.DeclineJoinRequestResponse{}, nil
}
