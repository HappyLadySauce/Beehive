package logic

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/model"
	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// RegisterLogic 实现用户注册：校验用户名与密码 → 检查用户名未被占用 → 写 users 表 → 生成 token 并返回（与登录一致，注册即登录）。
type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterRequest) (*pb.LoginResponse, error) {
	if in.GetUsername() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	_, err := l.svcCtx.UserMod.FindByUsername(in.GetUsername())
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "username already exists")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "hash password failed: %v", err)
	}

	// 10 位数字用户 ID（1000000000–9999999999），冲突则重试
	userID := generateTenDigitUserID()
	user := &model.User{
		ID:           userID,
		Username:     in.GetUsername(),
		PasswordHash: string(hash),
		Status:       "normal",
	}
	const maxRetries = 10
	created := false
	for i := 0; i < maxRetries; i++ {
		if err := l.svcCtx.UserMod.Create(user); err != nil {
			var pqErr *pq.Error
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				userID = generateTenDigitUserID()
				user.ID = userID
				continue
			}
			return nil, status.Errorf(codes.Internal, "create user failed: %v", err)
		}
		created = true
		break
	}
	if !created {
		return nil, status.Errorf(codes.Internal, "create user failed: too many id conflicts, please retry")
	}

	roles, err := l.svcCtx.RBACMod.GetUserRoles(userID)
	if err != nil {
		roles = nil
	}
	accessTTL := tokenTTLSeconds(l.svcCtx.Config.AccessTokenTTLSeconds, 3600)
	refreshTTL := tokenTTLSeconds(l.svcCtx.Config.RefreshTokenTTLSeconds, 30*24*3600)
	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()
	if err := storeToken(l.ctx, l.svcCtx.Redis, accessToken, userID, roles, time.Duration(accessTTL)*time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "store access token failed: %v", err)
	}
	if err := storeToken(l.ctx, l.svcCtx.Redis, refreshToken, userID, roles, time.Duration(refreshTTL)*time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "store refresh token failed: %v", err)
	}

	return &pb.LoginResponse{
		UserId:       userID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTTL),
	}, nil
}

// generateTenDigitUserID 生成 1000000000–9999999999 范围内的随机 10 位数字字符串
func generateTenDigitUserID() string {
	n := 1000000000 + rand.Int63n(9000000000)
	return strconv.FormatInt(n, 10)
}
