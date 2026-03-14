package logic

import (
	"context"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/auth/pb"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// LoginLogic 负责实现用户名密码登录的业务逻辑。
// 流程为：校验用户名与密码非空 → 查用户并校验状态与密码 → 拉取用户角色 → 生成 access/refresh token 并写入 Redis → 返回 token 与过期时间。
type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewLoginLogic 构造一个登录逻辑实例，注入 ServiceContext 以访问用户模型、RBAC、Redis 及配置。
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Login 使用用户名与密码完成登录。
// - 要求用户名、密码均非空；
// - 查用户并校验状态为 normal、密码通过 bcrypt 校验；
// - 拉取用户角色后生成 access 与 refresh token，分别写入 Redis 并设置 TTL；
// - 返回用户 ID、双 token 及 access token 的过期秒数。
func (l *LoginLogic) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 1. 参数校验：用户名、密码非空。
	if in.GetUsername() == "" || in.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "username and password are required")
	}

	// 2. 查用户并校验状态：
	//    - 用户不存在时与密码错误统一返回「用户名或密码错误」防枚举；
	//    - 仅当用户状态为 normal 时允许登录。
	user, err := l.svcCtx.UserMod.FindByUsername(in.GetUsername())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.Unauthenticated, "invalid username or password")
		}
		return nil, status.Errorf(codes.Internal, "query user failed: %v", err)
	}

	if user.Status != "normal" {
		return nil, status.Error(codes.FailedPrecondition, "account is not in normal status")
	}

	// 3. 密码校验：使用 bcrypt 比对请求密码与存储的哈希；失败时与用户不存在统一文案防枚举。
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.GetPassword())); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	// 4. 拉取用户角色，用于写入 token 载荷。
	roles, err := l.svcCtx.RBACMod.GetUserRoles(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get user roles failed: %v", err)
	}

	// 5. 从配置读取 TTL（未配置时使用默认值），生成 access/refresh token 并写入 Redis。
	accessTTL := tokenTTLSeconds(l.svcCtx.Config.AccessTokenTTLSeconds, 3600)
	refreshTTL := tokenTTLSeconds(l.svcCtx.Config.RefreshTokenTTLSeconds, 30*24*3600)

	accessToken := uuid.NewString()
	refreshToken := uuid.NewString()

	if err := storeToken(l.ctx, l.svcCtx.Redis, accessToken, user.ID, roles, time.Duration(accessTTL)*time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "store access token failed: %v", err)
	}
	if err := storeToken(l.ctx, l.svcCtx.Redis, refreshToken, user.ID, roles, time.Duration(refreshTTL)*time.Second); err != nil {
		return nil, status.Errorf(codes.Internal, "store refresh token failed: %v", err)
	}

	// 6. 返回用户 ID、双 token 及 access token 过期秒数。
	return &pb.LoginResponse{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTTL),
	}, nil
}
