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

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ctx := l.ctx
	id := in.GetId()
	key := fmt.Sprintf("user:profile:%s", id)

	// 1. 先查 Redis 缓存
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

	// 2. 缓存未命中，从 PostgreSQL 读取
	profile, err := l.svcCtx.UserProfileMod.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user profile failed: %v", err)
	}

	u := toProtoUser(profile)

	// 3. 回写缓存（错误只记日志，不影响主流程）
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
