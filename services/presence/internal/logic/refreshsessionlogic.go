package logic

import (
	"context"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/session"
	"github.com/HappyLadySauce/Beehive/services/presence/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RefreshSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshSessionLogic {
	return &RefreshSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefreshSessionLogic) RefreshSession(in *pb.RefreshSessionRequest) (*pb.RefreshSessionResponse, error) {
	if in.GetUserId() == "" || in.GetConnId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conn_id are required")
	}

	sessionKey := session.SessionKey(in.GetUserId(), in.GetConnId())
	ttl := session.SessionTTL(l.svcCtx.Config.SessionTTLSeconds)

	// 若 key 不存在（已过期）则视为幂等成功
	exists, err := l.svcCtx.Redis.Exists(l.ctx, sessionKey).Result()
	if err != nil {
		l.Errorf("refresh session exists check error: %v", err)
		return nil, status.Errorf(codes.Internal, "refresh session failed: %v", err)
	}
	if exists == 0 {
		return &pb.RefreshSessionResponse{}, nil
	}

	now := time.Now().Unix()
	if err := l.svcCtx.Redis.HSet(l.ctx, sessionKey, session.HashLastPingAt, strconv.FormatInt(now, 10)).Err(); err != nil {
		l.Errorf("refresh session hset error: %v", err)
		return nil, status.Errorf(codes.Internal, "refresh session failed: %v", err)
	}
	if err := l.svcCtx.Redis.Expire(l.ctx, sessionKey, time.Duration(ttl)*time.Second).Err(); err != nil {
		l.Errorf("refresh session expire error: %v", err)
		return nil, status.Errorf(codes.Internal, "refresh session failed: %v", err)
	}

	return &pb.RefreshSessionResponse{}, nil
}
