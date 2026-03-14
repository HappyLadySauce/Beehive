package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LogoutLogic 负责实现登出业务逻辑。
// 根据 access_token 删除 Redis 中对应的 token 会话，使该 token 失效。
type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewLogoutLogic 构造一个登出逻辑实例。
func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Logout 使当前 access_token 失效。
// - 要求 access_token 非空；
// - 校验 Redis 可用后，删除该 token 在 Redis 中的 key。
func (l *LogoutLogic) Logout(in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	// 1. 参数与 Redis 校验。
	if in.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	if l.svcCtx.Redis == nil {
		return nil, status.Error(codes.Internal, "redis client is not available")
	}

	// 2. 删除 token 对应 key，登出即删当前会话。
	if err := l.svcCtx.Redis.Del(l.ctx, tokenKey(in.GetAccessToken())).Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "logout failed: %v", err)
	}

	return &pb.LogoutResponse{}, nil
}
