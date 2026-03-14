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

type GetConversationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetConversationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetConversationLogic {
	return &GetConversationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetConversationLogic) GetConversation(in *pb.GetConversationRequest) (*pb.GetConversationResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	c, err := l.svcCtx.Conv.FindByID(in.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "conversation not found")
		}
		l.Errorf("find conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "find conversation failed: %v", err)
	}
	count, err := l.svcCtx.Conv.CountMembers(c.ID)
	if err != nil {
		l.Errorf("count members failed: %v", err)
		return nil, status.Errorf(codes.Internal, "count members failed: %v", err)
	}
	return &pb.GetConversationResponse{
		Conversation: &pb.ConversationInfo{
			Id:            c.ID,
			Type:          c.Type,
			Name:          c.Name,
			MemberCount:   int32(count),
			CreatedAt:     c.CreatedAt.Unix(),
			LastActiveAt:  c.LastActiveAt.Unix(),
		},
	}, nil
}
