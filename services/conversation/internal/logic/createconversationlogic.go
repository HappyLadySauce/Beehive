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
)

type CreateConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateConversationLogic {
	return &CreateConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateConversationLogic) CreateConversation(in *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	convType := in.GetType()
	if convType == "" {
		convType = "single"
	}
	now := time.Now()
	convID := uuid.Must(uuid.NewUUID()).String()
	conv := &model.Conversation{
		ID:           convID,
		Type:         convType,
		Name:         in.GetName(),
		CreatedAt:    now,
		LastActiveAt: now,
	}
	var members []*model.ConversationMember
	for i, uid := range in.GetMemberIds() {
		if uid == "" {
			continue
		}
		role := "member"
		if i == 0 {
			role = "owner"
		}
		members = append(members, &model.ConversationMember{
			ID:             uuid.Must(uuid.NewUUID()).String(),
			ConversationID: convID,
			UserID:         uid,
			Role:           role,
			Status:         "active",
			JoinedAt:       now,
		})
	}
	if err := l.svcCtx.Conv.Create(conv, members); err != nil {
		l.Errorf("create conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "create conversation failed: %v", err)
	}
	return &pb.CreateConversationResponse{ConversationId: convID}, nil
}