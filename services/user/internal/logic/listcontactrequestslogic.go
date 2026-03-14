package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListContactRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListContactRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListContactRequestsLogic {
	return &ListContactRequestsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListContactRequestsLogic) ListContactRequests(in *pb.ListContactRequestsRequest) (*pb.ListContactRequestsResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	list, err := l.svcCtx.ContactRequestMod.ListPendingByToUser(in.GetUserId())
	if err != nil {
		l.Errorf("list contact requests failed: %v", err)
		return nil, status.Errorf(codes.Internal, "list contact requests failed: %v", err)
	}
	items := make([]*pb.ContactRequestItem, 0, len(list))
	for _, r := range list {
		items = append(items, &pb.ContactRequestItem{
			RequestId:  r.ID,
			FromUserId: r.FromUserID,
			ToUserId:   r.ToUserID,
			Message:    r.Message,
			CreatedAt:  r.CreatedAt.Unix(),
		})
	}
	return &pb.ListContactRequestsResponse{Items: items}, nil
}
