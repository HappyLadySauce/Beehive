package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/HappyLadySauce/Beehive/services/user/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/user/pb"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BatchGetUsersLogic 负责实现批量查询用户资料的业务逻辑。
// 读取路径为：优先走 Redis 缓存（MGET），缓存未命中的用户再回源 PostgreSQL，
// 同时把从数据库查到的结果回写到 Redis，最后按照请求中的用户 ID 顺序构造返回列表。
type BatchGetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewBatchGetUsersLogic 构造一个批量查询用户逻辑对象。
// 这里把 go-zero 的日志上下文与业务上下文一起存入结构体，后续处理过程中都可以复用。
func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchGetUsers 批量获取用户基础资料。
// - 如果请求中没有任何 ID，则直接返回空列表；
// - 有 ID 时先做去重与顺序整理，再批量从 Redis 读取；
// - 缓存未命中的 ID 再从数据库中补齐，并把结果写回 Redis 做缓存。
func (l *BatchGetUsersLogic) BatchGetUsers(in *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	if len(in.GetIds()) == 0 {
		return &pb.BatchGetUsersResponse{Users: []*pb.User{}}, nil
	}

	// 去重并保持顺序：
	// - 使用 map 记录已经出现过的 ID，避免重复查询与返回；
	// - 遍历时依然按原始顺序把首次出现的 ID 追加到切片中，保证最终返回顺序与请求一致。
	seen := make(map[string]struct{}, len(in.Ids))
	var ids []string
	for _, id := range in.GetIds() {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return &pb.BatchGetUsersResponse{Users: []*pb.User{}}, nil
	}

	// ctx 统一从逻辑结构体中获取，保证链路上的日志与超时控制一致。
	ctx := l.ctx
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("user:profile:%s", id))
	}

	usersByID := make(map[string]*pb.User, len(ids))
	var missIDs []string

	// 1. 优先从 Redis 批量读取用户资料：
	//    使用 MGET 一次性拉取多个 key，尽可能降低网络往返次数。
	values, err := l.svcCtx.Redis.MGet(ctx, keys...).Result()
	if err != nil && err != redis.Nil {
		l.Errorf("redis MGET error: %v", err)
	}
	for i, v := range values {
		if v == nil {
			missIDs = append(missIDs, ids[i])
			continue
		}
		str, ok := v.(string)
		if !ok {
			missIDs = append(missIDs, ids[i])
			continue
		}
		var u pb.User
		if e := json.Unmarshal([]byte(str), &u); e != nil {
			missIDs = append(missIDs, ids[i])
			continue
		}
		usersByID[ids[i]] = &u
	}

	// 2. 对缓存未命中的 ID（missIDs）从 PostgreSQL 查询：
	//    - 只对缺失的数据访问数据库，避免对所有 ID 再查一遍；
	//    - 查到后不仅放入结果集，还会写入 Redis，后续请求可以直接命中缓存。
	if len(missIDs) > 0 {
		profiles, err := l.svcCtx.UserProfileMod.FindByIDs(missIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "batch query user profiles failed: %v", err)
		}
		if len(profiles) > 0 {
			// 构造缓存写入 pipeline：
			// - 使用 Redis Pipeline 可以把多次写操作打包发送，减少网络开销；
			// - TTL 从配置中读取，如果配置未设置则回退到 600 秒的默认值。
			pipe := l.svcCtx.Redis.Pipeline()
			ttl := l.svcCtx.Config.UserProfileTTLSeconds
			if ttl <= 0 {
				ttl = 600
			}
			exp := time.Duration(ttl) * time.Second

			for _, p := range profiles {
				// 将数据库模型转成对外的 protobuf User 结构，
				// 一方面用于组装返回值，另一方面用于写入缓存。
				u := toProtoUser(p)
				usersByID[p.UserID] = u
				if buf, e := json.Marshal(u); e == nil {
					key := fmt.Sprintf("user:profile:%s", p.UserID)
					pipe.Set(ctx, key, buf, exp)
				}
			}
			if _, e := pipe.Exec(ctx); e != nil {
				l.Errorf("redis pipeline SET profiles error: %v", e)
			}
		}
	}

	// 3. 按去重后的请求顺序组装返回，只返回成功获取到的用户：
	//    - 使用 ids（已经去重且保持顺序）驱动遍历；
	//    - 只有在 usersByID 中存在的 ID 才会被追加到结果列表中。
	result := make([]*pb.User, 0, len(ids))
	for _, id := range ids {
		if u, ok := usersByID[id]; ok {
			result = append(result, u)
		}
	}

	return &pb.BatchGetUsersResponse{Users: result}, nil
}
