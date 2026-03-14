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

type ApplyJoinGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyJoinGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyJoinGroupLogic {
	return &ApplyJoinGroupLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyJoinGroupLogic) ApplyJoinGroup(in *pb.ApplyJoinGroupRequest) (*pb.ApplyJoinGroupResponse, error) {
	if in.GetConversationId() == "" || in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id and user_id are required")
	}
	c, err := l.svcCtx.Conv.FindByID(in.GetConversationId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "conversation not found")
		}
		return nil, status.Errorf(codes.Internal, "find conversation failed: %v", err)
	}
	if c.Type != "group" {
		return nil, status.Error(codes.InvalidArgument, "not a group conversation")
	}
	member, _ := l.svcCtx.Conv.GetMember(in.GetConversationId(), in.GetUserId())
	if member != nil && member.Status == "active" {
		return nil, status.Error(codes.AlreadyExists, "already in group")
	}
	joinType := c.JoinType
	if joinType == "" {
		joinType = "approval"
	}
	if joinType == "direct" {
		now := time.Now()
		if err := l.svcCtx.Conv.AddMember(&model.ConversationMember{
			ID:             uuid.New().String(),
			ConversationID: in.GetConversationId(),
			UserID:         in.GetUserId(),
			Role:           "member",
			Status:         "active",
			JoinedAt:       now,
		}); err != nil {
			l.Errorf("add member for direct join failed: %v", err)
			return nil, status.Errorf(codes.Internal, "join failed: %v", err)
		}
		return &pb.ApplyJoinGroupResponse{Joined: true}, nil
	}
	req, err := l.svcCtx.JoinReq.Apply(in.GetConversationId(), in.GetUserId(), in.GetMessage())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.AlreadyExists, "already in group")
		}
		l.Errorf("apply join group failed: %v", err)
		return nil, status.Errorf(codes.Internal, "apply join group failed: %v", err)
	}
	return &pb.ApplyJoinGroupResponse{RequestId: req.ID, Joined: false}, nil
}
