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

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	ctx := l.ctx
	id := in.GetId()

	// 1. 先查当前 profile
	existing, err := l.svcCtx.UserProfileMod.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "query user profile failed: %v", err)
	}

	// 2. 根据请求选择性更新字段
	p := &model.UserProfile{
		UserID:    existing.UserID,
		Nickname:  chooseNonEmpty(in.GetNickname(), existing.Nickname),
		AvatarURL: chooseNonEmpty(in.GetAvatarUrl(), existing.AvatarURL),
		Bio:       chooseNonEmpty(in.GetBio(), existing.Bio),
		Status:    existing.Status,
		UpdatedAt: time.Now(),
	}

	if err := l.svcCtx.UserProfileMod.UpdateProfile(p); err != nil {
		return nil, status.Errorf(codes.Internal, "update user profile failed: %v", err)
	}

	u := toProtoUser(p)

	// 3. 刷新缓存
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

func chooseNonEmpty(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
