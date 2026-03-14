package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TokenLoginLogic 负责使用已有 access_token 做「免密登录」。
// 校验 token 有效后返回同一 token 及剩余过期时间，不签发新 token；RefreshToken 为空。
type TokenLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewTokenLoginLogic 构造一个 token 登录逻辑实例。
func NewTokenLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TokenLoginLogic {
	return &TokenLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// TokenLogin 用 access_token 完成登录态校验并返回登录结果；校验成功时会续期该 token 的 TTL，保持会话活跃。
// - 要求 access_token 非空；
// - 从 Redis 加载并校验 token，无效则返回错误；
// - 成功时对 key 执行 Expire 续期（与 Login 相同的 access TTL），再返回 LoginResponse：AccessToken 沿用入参、RefreshToken 为空、ExpiresIn 为续期后的 TTL 秒数。
func (l *TokenLoginLogic) TokenLogin(in *pb.TokenLoginRequest) (*pb.LoginResponse, error) {
	// 1. 参数校验。
	if in.GetAccessToken() == "" {
		return nil, status.Error(codes.InvalidArgument, "access_token is required")
	}

	// 2. 加载并校验 token，并续期 TTL（与 Login 使用相同 access TTL）；无效或不存在时返回 Unauthenticated，系统错误返回 Internal。
	accessTTL := tokenTTLSeconds(l.svcCtx.Config.AccessTokenTTLSeconds, 3600)
	extendTTL := time.Duration(accessTTL) * time.Second
	userID, _, ttl, err := loadAndTouchToken(l.ctx, l.svcCtx.Redis, in.GetAccessToken(), extendTTL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token login failed: %v", err)
	}
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "token invalid or expired")
	}

	// 3. 组装返回：同一 token + 续期后的 TTL，不返回 RefreshToken。
	return &pb.LoginResponse{
		UserId:       userID,
		AccessToken:  in.GetAccessToken(),
		RefreshToken: "",
		ExpiresIn:    int64(ttl.Seconds()),
	}, nil
}
