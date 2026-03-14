package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/model"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ApproveJoinRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApproveJoinRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApproveJoinRequestLogic {
	return &ApproveJoinRequestLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApproveJoinRequestLogic) ApproveJoinRequest(in *pb.ApproveJoinRequestRequest) (*pb.ApproveJoinRequestResponse, error) {
	if in.GetConversationId() == "" || in.GetRequestId() == "" || in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id, request_id and user_id are required")
	}
	member, err := l.svcCtx.Conv.GetMember(in.GetConversationId(), in.GetUserId())
	if err != nil || member == nil || member.Status != "active" {
		return nil, status.Error(codes.PermissionDenied, "not a member or not active")
	}
	if member.Role != "owner" && member.Role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "only owner or admin can approve")
	}
	req, err := l.svcCtx.JoinReq.FindByID(in.GetRequestId())
	if err != nil || req == nil || req.ConversationID != in.GetConversationId() || req.Status != "pending" {
		return nil, status.Error(codes.NotFound, "request not found or not pending")
	}
	if err := l.svcCtx.JoinReq.Approve(in.GetRequestId(), in.GetConversationId(), in.GetUserId()); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "request not found or not pending")
		}
		return nil, status.Errorf(codes.Internal, "approve failed: %v", err)
	}
	// 将申请人加入群
	now := time.Now()
	_ = l.svcCtx.Conv.AddMember(&model.ConversationMember{
		ID:             uuid.New().String(),
		ConversationID: in.GetConversationId(),
		UserID:         req.UserID,
		Role:           "member",
		Status:         "active",
		JoinedAt:       now,
	})
	return &pb.ApproveJoinRequestResponse{}, nil
}
