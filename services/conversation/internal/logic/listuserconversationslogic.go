package logic

import (
	"context"
	"strconv"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListUserConversationsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserConversationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserConversationsLogic {
	return &ListUserConversationsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListUserConversationsLogic) ListUserConversations(in *pb.ListUserConversationsRequest) (*pb.ListUserConversationsResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	offset := 0
	if in.GetCursor() != "" {
		if o, err := strconv.Atoi(in.GetCursor()); err == nil && o >= 0 {
			offset = o
		}
	}
	limit := int(in.GetLimit())
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	list, err := l.svcCtx.Conv.ListByUserID(in.GetUserId(), offset, limit+1)
	if err != nil {
		l.Errorf("list user conversations failed: %v", err)
		return nil, status.Errorf(codes.Internal, "list user conversations failed: %v", err)
	}
	hasMore := len(list) > limit
	if hasMore {
		list = list[:limit]
	}
	items := make([]*pb.ConversationInfo, 0, len(list))
	for _, c := range list {
		count, _ := l.svcCtx.Conv.CountMembers(c.ID)
		items = append(items, &pb.ConversationInfo{
			Id:            c.ID,
			Type:          c.Type,
			Name:          c.Name,
			MemberCount:   int32(count),
			CreatedAt:     c.CreatedAt.Unix(),
			LastActiveAt:  c.LastActiveAt.Unix(),
		})
	}
	var nextCursor string
	if hasMore {
		nextCursor = strconv.Itoa(offset + limit)
	}
	return &pb.ListUserConversationsResponse{Items: items, NextCursor: nextCursor}, nil
}
