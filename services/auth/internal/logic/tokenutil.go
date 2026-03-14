package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// tokenPayload 为 Redis 中 token 对应的 JSON 结构，用于存储 user_id 与 roles。
type tokenPayload struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
}

// tokenTTLSeconds 用于 access/refresh token 的 TTL 取值：配置值 ≤0 时使用默认值 def，否则使用配置值。
func tokenTTLSeconds(v int, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

// tokenKey 返回统一的 Redis key 前缀 "auth:token:" + token。
func tokenKey(token string) string {
	return "auth:token:" + token
}

// storeToken 将 userID、roles 序列化为 JSON 写入 Redis，并设置 ttl；rdb 为 nil 时返回错误。
func storeToken(ctx context.Context, rdb *redis.Client, token string, userID string, roles []string, ttl time.Duration) error {
	if rdb == nil {
		return errors.New("redis client is nil")
	}
	payload := tokenPayload{
		UserID: userID,
		Roles:  roles,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal token payload failed: %w", err)
	}
	return rdb.Set(ctx, tokenKey(token), b, ttl).Err()
}

// loadToken 从 Redis 读取 token 并解析出 userID、roles，同时返回剩余 TTL。
// key 不存在或为 redis.Nil 时返回空 userID、nil roles、0、nil error；不做续期。
func loadToken(ctx context.Context, rdb *redis.Client, token string) (string, []string, time.Duration, error) {
	if rdb == nil {
		return "", nil, 0, errors.New("redis client is nil")
	}
	key := tokenKey(token)
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil, 0, nil
		}
		return "", nil, 0, fmt.Errorf("redis get failed: %w", err)
	}

	var payload tokenPayload
	if err := json.Unmarshal([]byte(res), &payload); err != nil {
		return "", nil, 0, fmt.Errorf("unmarshal token payload failed: %w", err)
	}

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return payload.UserID, payload.Roles, 0, fmt.Errorf("redis ttl failed: %w", err)
	}

	return payload.UserID, payload.Roles, ttl, nil
}

