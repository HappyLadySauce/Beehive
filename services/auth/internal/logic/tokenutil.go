package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type tokenPayload struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
}

func tokenTTLSeconds(v int, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func tokenKey(token string) string {
	return "auth:token:" + token
}

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

func loadAndTouchToken(ctx context.Context, rdb *redis.Client, token string) (string, []string, time.Duration, error) {
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

