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

type BatchGetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUsersLogic) BatchGetUsers(in *pb.BatchGetUsersRequest) (*pb.BatchGetUsersResponse, error) {
	if len(in.GetIds()) == 0 {
		return &pb.BatchGetUsersResponse{Users: []*pb.User{}}, nil
	}

	// 去重并保持顺序
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

	ctx := l.ctx
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, fmt.Sprintf("user:profile:%s", id))
	}

	usersByID := make(map[string]*pb.User, len(ids))
	var missIDs []string

	// 1. Redis MGet
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

	// 2. 对 miss 的 ID 从 PostgreSQL 查询
	if len(missIDs) > 0 {
		profiles, err := l.svcCtx.UserProfileMod.FindByIDs(missIDs)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "batch query user profiles failed: %v", err)
		}
		if len(profiles) > 0 {
			// 构造缓存写入 pipeline
			pipe := l.svcCtx.Redis.Pipeline()
			ttl := l.svcCtx.Config.UserProfileTTLSeconds
			if ttl <= 0 {
				ttl = 600
			}
			exp := time.Duration(ttl) * time.Second

			for _, p := range profiles {
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

	// 3. 按去重后的请求顺序组装返回，只返回存在的用户
	result := make([]*pb.User, 0, len(ids))
	for _, id := range ids {
		if u, ok := usersByID[id]; ok {
			result = append(result, u)
		}
	}

	return &pb.BatchGetUsersResponse{Users: result}, nil
}
