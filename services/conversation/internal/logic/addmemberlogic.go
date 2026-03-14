package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/model"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AddMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMemberLogic {
	return &AddMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMemberLogic) AddMember(in *pb.AddMemberRequest) (*pb.AddMemberResponse, error) {
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
	role := in.GetRole()
	if role == "" {
		role = "member"
	}
	_, err = l.svcCtx.Conv.GetMember(in.GetConversationId(), in.GetUserId())
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "user already in conversation")
	}
	if err != gorm.ErrRecordNotFound {
		l.Errorf("get member failed: %v", err)
		return nil, status.Errorf(codes.Internal, "get member failed: %v", err)
	}
	member := &model.ConversationMember{
		ID:             uuid.Must(uuid.NewUUID()).String(),
		ConversationID: in.GetConversationId(),
		UserID:         in.GetUserId(),
		Role:           role,
		Status:         "active",
	}
	if err := l.svcCtx.Conv.AddMember(member); err != nil {
		l.Errorf("add member failed: %v", err)
		return nil, status.Errorf(codes.Internal, "add member failed: %v", err)
	}
	return &pb.AddMemberResponse{}, nil
}
