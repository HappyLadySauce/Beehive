package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HappyLadySauce/Beehive/services/user/internal/model"
	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

// UpdateUserLogic 负责实现用户资料更新相关的业务逻辑。
// 更新流程大致为：
// 1. 先从数据库中读取当前用户 profile；
// 2. 根据请求参数有选择地覆盖可变字段（昵称、头像、个性签名等）；
// 3. 更新数据库记录；
// 4. 将最新数据写回 Redis 缓存，保证后续读操作可以直接命中缓存。
type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewUpdateUserLogic 构造一个用户更新逻辑实例。
// 这里注入的 ServiceContext 提供了数据库模型、Redis 客户端以及配置等依赖。
func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUser 更新用户的基础资料。
// - 要求请求中必须携带用户 ID；
// - 仅对非空字段做覆盖更新，保持“部分更新”语义；
// - 更新数据库成功后会同步刷新 Redis 缓存中的用户资料。
func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ctx := l.ctx
	id := in.GetId()

	// 1. 先查当前 profile：
	//    - 如果记录不存在，返回 NotFound；
	//    - 其他数据库错误统一包装为 Internal。
	existing, err := l.svcCtx.UserProfileMod.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user profile failed: %v", err)
	}

	// 2. 根据请求选择性更新字段：
	//    - 使用 chooseNonEmpty 对字符串字段进行“非空覆盖”，
	//      即只有当请求字段非空时才会覆盖原有值，避免把空字符串写回数据库。
	p := &model.UserProfile{
		UserID:    existing.UserID,
		Nickname:  chooseNonEmpty(in.GetNickname(), existing.Nickname),
		AvatarURL: chooseNonEmpty(in.GetAvatarUrl(), existing.AvatarURL),
		Bio:       chooseNonEmpty(in.GetBio(), existing.Bio),
		Status:    existing.Status,
		UpdatedAt: time.Now(),
	}

	// 3. 写回数据库，持久化更新后的用户资料。
	if err := l.svcCtx.UserProfileMod.UpdateProfile(p); err != nil {
		return nil, status.Errorf(codes.Internal, "update user profile failed: %v", err)
	}

	// 将内部模型转换为对外返回的 protobuf User 结构。
	u := toProtoUser(p)

	// 4. 刷新缓存：
	//    - 把最新的用户资料序列化为 JSON 存入 Redis；
	//    - TTL 同样来自配置，未设置时默认使用 600 秒；
	//    - 如果写缓存失败，只记录日志，不影响更新流程的主路径。
	key := fmt.Sprintf("user:profile:%s", id)
	if buf, e := json.Marshal(u); e == nil {
		ttl := l.svcCtx.Config.UserProfileTTLSeconds
		if ttl <= 0 {
			ttl = 600
		}
		if e := l.svcCtx.Redis.Set(ctx, key, buf, time.Duration(ttl)*time.Second).Err(); e != nil && e != redis.Nil {
			l.Errorf("redis SET %s error: %v", key, e)
		}
	} else {
		l.Errorf("marshal updated user profile failed: %v", e)
	}

	return &pb.UpdateUserResponse{User: u}, nil
}

// chooseNonEmpty 在更新逻辑中实现“非空覆盖”策略：
// - 如果传入的新值 v 非空，则使用 v；
// - 否则回退到原有值 fallback。
func chooseNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
