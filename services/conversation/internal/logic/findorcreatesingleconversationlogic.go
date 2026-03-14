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

type FindOrCreateSingleConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindOrCreateSingleConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindOrCreateSingleConversationLogic {
	return &FindOrCreateSingleConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindOrCreateSingleConversationLogic) FindOrCreateSingleConversation(in *pb.FindOrCreateSingleConversationRequest) (*pb.FindOrCreateSingleConversationResponse, error) {
	userID1 := in.GetUserId_1()
	userID2 := in.GetUserId_2()
	if userID1 == "" || userID2 == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id_1 and user_id_2 are required")
	}
	if userID1 == userID2 {
		return nil, status.Error(codes.InvalidArgument, "user_id_1 and user_id_2 must be different")
	}
	conv, err := l.svcCtx.Conv.FindSingleByTwoUsers(userID1, userID2)
	if err == nil {
		return &pb.FindOrCreateSingleConversationResponse{ConversationId: conv.ID}, nil
	}
	if err != gorm.ErrRecordNotFound {
		l.Errorf("find single conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "find single conversation failed: %v", err)
	}
	// 不存在则创建单聊
	now := time.Now()
	convID := uuid.Must(uuid.NewUUID()).String()
	newConv := &model.Conversation{
		ID:           convID,
		Type:         "single",
		Name:         "",
		CreatedAt:    now,
		LastActiveAt: now,
	}
	members := []*model.ConversationMember{
		{ID: uuid.Must(uuid.NewUUID()).String(), ConversationID: convID, UserID: userID1, Role: "owner", Status: "active", JoinedAt: now},
		{ID: uuid.Must(uuid.NewUUID()).String(), ConversationID: convID, UserID: userID2, Role: "member", Status: "active", JoinedAt: now},
	}
	if err = l.svcCtx.Conv.Create(newConv, members); err != nil {
		l.Errorf("create single conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "create single conversation failed: %v", err)
	}
	return &pb.FindOrCreateSingleConversationResponse{ConversationId: convID}, nil
}
