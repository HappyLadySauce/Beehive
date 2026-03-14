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

type ListMembersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMembersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMembersLogic {
	return &ListMembersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMembersLogic) ListMembers(in *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	if in.GetConversationId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id is required")
	}
	_, err := l.svcCtx.Conv.FindByID(in.GetConversationId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "conversation not found")
		}
		l.Errorf("find conversation failed: %v", err)
		return nil, status.Errorf(codes.Internal, "find conversation failed: %v", err)
	}
	list, err := l.svcCtx.Conv.ListMembers(in.GetConversationId())
	if err != nil {
		l.Errorf("list members failed: %v", err)
		return nil, status.Errorf(codes.Internal, "list members failed: %v", err)
	}
	items := make([]*pb.MemberInfo, 0, len(list))
	for _, m := range list {
		items = append(items, &pb.MemberInfo{
			UserId:   m.UserID,
			Role:     m.Role,
			JoinedAt: m.JoinedAt.Unix(),
			Status:   m.Status,
		})
	}
	return &pb.ListMembersResponse{Items: items}, nil
}
