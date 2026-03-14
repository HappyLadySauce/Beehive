package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/session"
	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UnregisterSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnregisterSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnregisterSessionLogic {
	return &UnregisterSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnregisterSessionLogic) UnregisterSession(in *pb.UnregisterSessionRequest) (*pb.UnregisterSessionResponse, error) {
	if in.GetUserId() == "" || in.GetConnId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conn_id are required")
	}

	userConnsKey := session.UserConnsKey(in.GetUserId())
	sessionKey := session.SessionKey(in.GetUserId(), in.GetConnId())

	pipe := l.svcCtx.Redis.Pipeline()
	pipe.SRem(l.ctx, userConnsKey, in.GetConnId())
	pipe.Del(l.ctx, sessionKey)
	_, err := pipe.Exec(l.ctx)
	if err != nil {
		l.Errorf("unregister session redis error: %v", err)
		return nil, status.Errorf(codes.Internal, "unregister session failed: %v", err)
	}

	// 若该用户已无任何会话，删除空 set，避免残留 key
	if card, e := l.svcCtx.Redis.SCard(l.ctx, userConnsKey).Result(); e == nil && card == 0 {
		_ = l.svcCtx.Redis.Del(l.ctx, userConnsKey).Err()
	}

	return &pb.UnregisterSessionResponse{}, nil
}
