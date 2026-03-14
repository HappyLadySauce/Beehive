package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/conversation/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/conversation/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListJoinRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListJoinRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListJoinRequestsLogic {
	return &ListJoinRequestsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListJoinRequestsLogic) ListJoinRequests(in *pb.ListJoinRequestsRequest) (*pb.ListJoinRequestsResponse, error) {
	if in.GetConversationId() == "" || in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "conversation_id and user_id are required")
	}
	member, err := l.svcCtx.Conv.GetMember(in.GetConversationId(), in.GetUserId())
	if err != nil || member == nil || member.Status != "active" {
		return nil, status.Error(codes.PermissionDenied, "not a member or not active")
	}
	if member.Role != "owner" && member.Role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "only owner or admin can list join requests")
	}
	list, err := l.svcCtx.JoinReq.ListPending(in.GetConversationId())
	if err != nil {
		l.Errorf("list join requests failed: %v", err)
		return nil, status.Errorf(codes.Internal, "list join requests failed: %v", err)
	}
	items := make([]*pb.JoinRequestItem, 0, len(list))
	for _, r := range list {
		items = append(items, &pb.JoinRequestItem{
			RequestId:      r.ID,
			ConversationId: r.ConversationID,
			UserId:         r.UserID,
			Message:        r.Message,
			CreatedAt:      r.CreatedAt.Unix(),
		})
	}
	return &pb.ListJoinRequestsResponse{Items: items}, nil
}
