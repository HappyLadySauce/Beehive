package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidateTokenLogic 负责校验 access_token 是否有效并返回对应用户 ID。
// 空 token 或 token 不存在/无效时返回 valid=false，不返回错误；仅当 Redis 等系统错误时返回 error。
type ValidateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewValidateTokenLogic 构造一个 token 校验逻辑实例。
func NewValidateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateTokenLogic {
	return &ValidateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ValidateToken 校验 access_token 是否有效。
// - 空 token 直接返回 valid=false、空 UserId，不报错；
// - 从 Redis 加载并解析 token，若不存在或解析后无 userID 则返回 valid=false；
// - 有效时返回 valid=true 及 UserId。
func (l *ValidateTokenLogic) ValidateToken(in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	// 1. 空 token 处理：直接返回 valid=false，不返回错误。
	if in.GetAccessToken() == "" {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
		}, nil
	}

	// 2. 从 Redis 加载并解析 token；若 load 失败（如 Redis 错误）则返回 Internal。
	userID, _, _, err := loadToken(l.ctx, l.svcCtx.Redis, in.GetAccessToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "validate token failed: %v", err)
	}
	if userID == "" {
		return &pb.ValidateTokenResponse{
			Valid:  false,
			UserId: "",
		}, nil
	}

	// 3. 有效 token，返回 valid=true 及 UserId。
	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
	}, nil
}
