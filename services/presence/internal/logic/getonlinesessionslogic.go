package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetOnlineSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOnlineSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineSessionsLogic {
	return &GetOnlineSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOnlineSessionsLogic) GetOnlineSessions(in *pb.GetOnlineSessionsRequest) (*pb.GetOnlineSessionsResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	sessions, err := getSessionsForUser(l.ctx, l.svcCtx.Redis, in.GetUserId())
	if err != nil {
		l.Errorf("get online sessions error: %v", err)
		return nil, status.Errorf(codes.Internal, "get online sessions failed: %v", err)
	}
	return &pb.GetOnlineSessionsResponse{Sessions: sessions}, nil
}
