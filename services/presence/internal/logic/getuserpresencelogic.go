package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetUserPresenceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserPresenceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPresenceLogic {
	return &GetUserPresenceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserPresenceLogic) GetUserPresence(in *pb.GetUserPresenceRequest) (*pb.GetUserPresenceResponse, error) {
	if in.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	sessions, err := getSessionsForUser(l.ctx, l.svcCtx.Redis, in.GetUserId())
	if err != nil {
		l.Errorf("get user presence error: %v", err)
		return nil, status.Errorf(codes.Internal, "get user presence failed: %v", err)
	}
	return &pb.GetUserPresenceResponse{
		Online:   len(sessions) > 0,
		Sessions: sessions,
	}, nil
}
