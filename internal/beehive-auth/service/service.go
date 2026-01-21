package service

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"

	pb "github.com/HappyLadySauce/Beehive/api/proto/auth/v1"
	userpb "github.com/HappyLadySauce/Beehive/api/proto/user/v1"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/client"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/store"
	"github.com/HappyLadySauce/Beehive/pkg/utils/jwt"
	"github.com/HappyLadySauce/Beehive/pkg/utils/passwd"
)

// Service Auth Service 实现
type Service struct {
	pb.UnimplementedAuthServiceServer
	config     *config.Config
	redisStore *store.Store
	userClient *client.Client
}

// NewService 创建新的 Auth Service
func NewService(cfg *config.Config, redisStore *store.Store, userClient *client.Client) *Service {
	return &Service{
		config:     cfg,
		redisStore: redisStore,
		userClient: userClient,
	}
}

// Login 用户登录
func (s *Service) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. 验证请求参数
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// 2. 调用 User Service 获取用户信息（包含密码哈希和盐值）
	userResp, err := s.userClient.GetUserByID(ctx, req.Id)
	if err != nil {
		klog.Warningf("Login failed: user not found, id=%s, error=%v", req.Id, err)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// 3. 验证密码
	if !passwd.VerifyPassword(req.Password, userResp.Salt, userResp.PasswordHash) {
		klog.Warningf("Login failed: invalid password, id=%s", req.Id)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// 4. 检查用户状态（是否被冻结）
	// 注意：UserInfo 中没有 FreezeTime 字段，这里假设如果用户存在且密码正确就可以登录
	// 如果需要检查冻结状态，需要从 User Service 返回更多信息

	// 5. 生成 Access Token
	accessToken, err := s.generateAccessToken(ctx, userResp.User.Id, userResp.User.Email)
	if err != nil {
		klog.Errorf("Failed to generate access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	// 6. 生成 Refresh Token
	refreshToken, err := s.generateRefreshToken(ctx, userResp.User.Id)
	if err != nil {
		klog.Errorf("Failed to generate refresh token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// 7. 计算过期时间
	expiresAt := time.Now().Add(time.Duration(s.config.JWT.ExpireHours) * time.Hour).Unix()

	klog.Infof("User logged in successfully: id=%s, email=%s", userResp.User.Id, userResp.User.Email)

	return &pb.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         convertUserInfo(userResp.User),
	}, nil
}

// ValidateToken 验证 JWT Token
func (s *Service) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// 1. 检查黑名单
	inBlacklist, err := s.redisStore.IsInBlacklist(ctx, req.Token)
	if err != nil {
		klog.Errorf("Failed to check blacklist: %v", err)
		// 继续验证，不因为 Redis 错误而拒绝
	}
	if inBlacklist {
		klog.V(4).Infof("Token is in blacklist: %s", req.Token)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// 2. 验证 JWT 签名和过期时间
	claims, err := jwt.ParseToken(req.Token, s.config.JWT.Secret)
	if err != nil {
		klog.V(4).Infof("Token validation failed: %v", err)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// 3. 检查 Token 类型（应该是 access token）
	// 注意：当前的 JWT Claims 结构中没有 type 字段，这里假设所有有效的 access token 都可以通过

	// 4. 检查缓存
	cachedUserID, err := s.redisStore.GetTokenCache(ctx, req.Token)
	if err == nil && cachedUserID != "" {
		// 缓存命中，查询用户信息
		userResp, err := s.userClient.GetUser(ctx, cachedUserID)
		if err == nil {
			return &pb.ValidateTokenResponse{
				Valid: true,
				Id:    claims.UserID,
				User:  convertUserInfo(userResp.User),
			}, nil
		}
	}

	// 5. 缓存未命中，查询用户信息
	userResp, err := s.userClient.GetUser(ctx, claims.UserID)
	if err != nil {
		klog.Warningf("User not found during token validation: id=%s, error=%v", claims.UserID, err)
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	// 6. 缓存结果
	tokenExpiration := time.Duration(s.config.JWT.ExpireHours) * time.Hour
	if err := s.redisStore.SetTokenCache(ctx, req.Token, claims.UserID, tokenExpiration); err != nil {
		klog.Warningf("Failed to cache token: %v", err)
		// 继续返回，不因为缓存失败而拒绝
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		Id:    claims.UserID,
		User:  convertUserInfo(userResp.User),
	}, nil
}

// RefreshToken 刷新访问令牌
func (s *Service) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "refresh token is required")
	}

	// 1. 从 Redis 获取 Refresh Token 对应的用户ID
	userID, err := s.redisStore.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		klog.Warningf("Invalid refresh token: %v", err)
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	// 2. 获取用户信息以获取 email
	userResp, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		klog.Errorf("Failed to get user info: %v", err)
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	// 3. 生成新的 Access Token
	newAccessToken, err := s.generateAccessToken(ctx, userID, userResp.User.Email)
	if err != nil {
		klog.Errorf("Failed to generate new access token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate token")
	}

	// 4. 生成新的 Refresh Token
	newRefreshToken, err := s.generateRefreshToken(ctx, userID)
	if err != nil {
		klog.Errorf("Failed to generate new refresh token: %v", err)
		return nil, status.Error(codes.Internal, "failed to generate refresh token")
	}

	// 5. 删除旧的 Refresh Token
	if err := s.redisStore.DeleteRefreshToken(ctx, req.RefreshToken); err != nil {
		klog.Warningf("Failed to delete old refresh token: %v", err)
		// 继续执行，不因为删除失败而拒绝
	}

	// 6. 计算过期时间
	expiresAt := time.Now().Add(time.Duration(s.config.JWT.ExpireHours) * time.Hour).Unix()

	klog.Infof("Token refreshed successfully: user_id=%s", userID)

	return &pb.RefreshTokenResponse{
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RevokeToken 撤销令牌（加入黑名单）
func (s *Service) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	// 1. 验证 Token 有效性（可选，即使无效也可以加入黑名单）
	claims, err := jwt.ParseToken(req.Token, s.config.JWT.Secret)
	var expiration time.Duration
	if err == nil && claims != nil {
		// Token 有效，使用其剩余过期时间
		if claims.ExpiresAt != nil {
			expTime := claims.ExpiresAt.Time
			expiration = time.Until(expTime)
			if expiration < 0 {
				expiration = time.Hour // 如果已过期，设置1小时
			}
		} else {
			expiration = time.Duration(s.config.JWT.ExpireHours) * time.Hour
		}
	} else {
		// Token 无效，使用默认过期时间
		expiration = time.Duration(s.config.JWT.ExpireHours) * time.Hour
	}

	// 2. 加入黑名单
	if err := s.redisStore.AddToBlacklist(ctx, req.Token, expiration); err != nil {
		klog.Errorf("Failed to add token to blacklist: %v", err)
		return nil, status.Error(codes.Internal, "failed to revoke token")
	}

	// 3. 清除缓存
	if err := s.redisStore.DeleteTokenCache(ctx, req.Token); err != nil {
		klog.Warningf("Failed to delete token cache: %v", err)
		// 继续执行，不因为删除失败而拒绝
	}

	klog.Infof("Token revoked successfully")

	return &pb.RevokeTokenResponse{
		Success: true,
	}, nil
}

// generateAccessToken 生成 Access Token
func (s *Service) generateAccessToken(ctx context.Context, userID, email string) (string, error) {
	expiration := time.Duration(s.config.JWT.ExpireHours) * time.Hour
	token, err := jwt.GenerateToken(userID, email, "user", s.config.JWT.Secret, expiration)
	if err != nil {
		return "", err
	}

	// 缓存 Token
	if err := s.redisStore.SetTokenCache(ctx, token, userID, expiration); err != nil {
		klog.Warningf("Failed to cache access token: %v", err)
		// 继续返回 token，不因为缓存失败而拒绝
	}

	return token, nil
}

// generateRefreshToken 生成 Refresh Token
func (s *Service) generateRefreshToken(ctx context.Context, userID string) (string, error) {
	// Refresh Token 也使用 JWT，但存储到 Redis
	expiration := time.Duration(s.config.JWT.RefreshExpireHours) * time.Hour
	// 使用空字符串作为 username 和 role，因为 Refresh Token 只需要 userID
	token, err := jwt.GenerateToken(userID, "", "refresh", s.config.JWT.Secret, expiration)
	if err != nil {
		return "", err
	}

	// 存储 Refresh Token 到 Redis
	if err := s.redisStore.SetRefreshToken(ctx, token, userID, expiration); err != nil {
		klog.Warningf("Failed to store refresh token: %v", err)
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}

// convertUserInfo 将 User Service 的 UserInfo 转换为 Auth Service 的 UserInfo
func convertUserInfo(user *userpb.UserInfo) *pb.UserInfo {
	if user == nil {
		return nil
	}
	return &pb.UserInfo{
		Id:          user.Id,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Description: user.Description,
		Level:       user.Level,
		Status:      user.Status,
	}
}
