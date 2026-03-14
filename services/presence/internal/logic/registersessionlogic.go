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

type RegisterSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterSessionLogic {
	return &RegisterSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterSessionLogic) RegisterSession(in *pb.RegisterSessionRequest) (*pb.RegisterSessionResponse, error) {
	if in.GetUserId() == "" || in.GetConnId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and conn_id are required")
	}

	ttl := session.SessionTTL(l.svcCtx.Config.SessionTTLSeconds)
	now := time.Now().Unix()

	userConnsKey := session.UserConnsKey(in.GetUserId())
	sessionKey := session.SessionKey(in.GetUserId(), in.GetConnId())

	pipe := l.svcCtx.Redis.Pipeline()
	pipe.SAdd(l.ctx, userConnsKey, in.GetConnId())
	pipe.HSet(l.ctx, sessionKey,
		session.HashGatewayID, in.GetGatewayId(),
		session.HashConnID, in.GetConnId(),
		session.HashDeviceID, in.GetDeviceId(),
		session.HashDeviceType, in.GetDeviceType(),
		session.HashLastPingAt, strconv.FormatInt(now, 10),
	)
	pipe.Expire(l.ctx, sessionKey, time.Duration(ttl)*time.Second)
	_, err := pipe.Exec(l.ctx)
	if err != nil {
		l.Errorf("register session redis error: %v", err)
		return nil, status.Errorf(codes.Internal, "register session failed: %v", err)
	}

	return &pb.RegisterSessionResponse{}, nil
}