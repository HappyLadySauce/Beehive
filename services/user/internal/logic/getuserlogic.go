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

// GetUserLogic 负责实现单个用户资料查询的业务逻辑。
// 数据读取路径为：Redis 缓存优先，缓存未命中或反序列化失败时回源 PostgreSQL，
// 同时在成功从数据库读到数据后，将结果序列化回写到 Redis 中。
type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewGetUserLogic 构造一个用户查询逻辑实例。
// 通过注入 ServiceContext，可以在逻辑层中访问 Redis、数据库模型等基础设施组件。
func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUser 查询单个用户的公开资料。
// - 入参要求提供用户 ID，ID 为空时直接返回 InvalidArgument 错误；
// - 先尝试从 Redis 读取并反序列化，如果成功则直接返回；
// - 如果缓存不存在或内容异常，再从 PostgreSQL 查询，并把查询结果写入缓存。
func (l *GetUserLogic) GetUser(in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ctx := l.ctx
	id := in.GetId()
	key := fmt.Sprintf("user:profile:%s", id)

	// 1. 先查 Redis 缓存：
	//    - key 采用统一前缀 "user:profile:" + 用户 ID；
	//    - Get 之后直接以字节切片方式取值，方便后续 JSON 反序列化。
	val, err := l.svcCtx.Redis.Get(ctx, key).Bytes()
	if err == nil && len(val) > 0 {
		var u pb.User
		if e := json.Unmarshal(val, &u); e == nil {
			return &pb.GetUserResponse{User: &u}, nil
		}
	}
	if err != nil && err != redis.Nil {
		l.Errorf("redis GET %s error: %v", key, err)
	}

	// 2. 缓存未命中或读取失败时，从 PostgreSQL 读取用户资料：
	//    - 若记录不存在，返回 NotFound；
	//    - 其他错误（例如数据库异常）统一包装为 Internal 错误。
	profile, err := l.svcCtx.UserProfileMod.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user profile failed: %v", err)
	}

	// 将数据库模型转换为对外返回的 protobuf User 结构体。
	u := toProtoUser(profile)

	// 3. 回写缓存（错误只记日志，不影响主流程）：
	//    - 序列化成功才会尝试写入 Redis；
	//    - TTL 由配置控制，未配置时采用 600 秒的默认值。
	if buf, e := json.Marshal(u); e == nil {
		ttl := l.svcCtx.Config.UserProfileTTLSeconds
		if ttl <= 0 {
			ttl = 600
		}
		if e := l.svcCtx.Redis.Set(ctx, key, buf, time.Duration(ttl)*time.Second).Err(); e != nil {
			l.Errorf("redis SET %s error: %v", key, e)
		}
	} else {
		l.Errorf("marshal user profile for cache failed: %v", e)
	}

	return &pb.GetUserResponse{User: u}, nil
}

func toProtoUser(p *model.UserProfile) *pb.User {
	return &pb.User{
		Id:        p.UserID,
		Nickname:  p.Nickname,
		AvatarUrl: p.AvatarURL,
		Bio:       p.Bio,
		Status:    p.Status,
	}
}
