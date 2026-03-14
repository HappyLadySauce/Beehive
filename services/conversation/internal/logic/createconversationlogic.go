package logic

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/model"
	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/google/uuid"
	"github.com/lib/pq"
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
	var convID string
	if convType == "group" {
		const maxRetries = 10
		for attempt := 0; attempt < maxRetries; attempt++ {
			convID = generateElevenDigitGroupID()
			conv := &model.Conversation{
				ID:           convID,
				Type:         convType,
				Name:         in.GetName(),
				CreatedAt:    now,
				LastActiveAt: now,
			}
			var members []*model.ConversationMember
			firstMember := true
			for _, uid := range in.GetMemberIds() {
				if uid == "" {
					continue
				}
				role := "member"
				if firstMember {
					role = "owner"
					firstMember = false
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
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code == "23505" {
					continue
				}
				l.Errorf("create conversation failed: %v", err)
				return nil, status.Errorf(codes.Internal, "create conversation failed: %v", err)
			}
			return &pb.CreateConversationResponse{ConversationId: convID}, nil
		}
		l.Errorf("create group conversation: too many id conflicts")
		return nil, status.Errorf(codes.Internal, "create conversation failed: id conflict")
	}
	convID = uuid.Must(uuid.NewUUID()).String()
	conv := &model.Conversation{
		ID:           convID,
		Type:         convType,
		Name:         in.GetName(),
		CreatedAt:    now,
		LastActiveAt: now,
	}
	var members []*model.ConversationMember
	firstMember := true
	for _, uid := range in.GetMemberIds() {
		if uid == "" {
			continue
		}
		role := "member"
		if firstMember {
			role = "owner"
			firstMember = false
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

// generateElevenDigitGroupID 生成 10000000000–99999999999 范围内的随机 11 位数字字符串（群聊 ID）
func generateElevenDigitGroupID() string {
	n := 10000000000 + rand.Int63n(90000000000)
	return strconv.FormatInt(n, 10)
}