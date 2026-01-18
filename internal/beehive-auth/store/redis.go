package store

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"k8s.io/klog/v2"

	"github.com/HappyLadySauce/Beehive/internal/beehive-auth/config"
)

// Store Redis 存储封装
type Store struct {
	client *redis.Client
}

// NewStore 创建新的 Redis 存储
func NewStore(cfg *config.Config) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	klog.Info("Redis connection established successfully")
	return &Store{
		client: client,
	}, nil
}

// Close 关闭 Redis 连接
func (s *Store) Close() error {
	return s.client.Close()
}

// Client 返回 Redis 客户端
func (s *Store) Client() *redis.Client {
	return s.client
}

// SetTokenCache 设置 Token 缓存
func (s *Store) SetTokenCache(ctx context.Context, token, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("token:%s", token)
	return s.client.Set(ctx, key, userID, expiration).Err()
}

// GetTokenCache 获取 Token 缓存
func (s *Store) GetTokenCache(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("token:%s", token)
	return s.client.Get(ctx, key).Result()
}

// DeleteTokenCache 删除 Token 缓存
func (s *Store) DeleteTokenCache(ctx context.Context, token string) error {
	key := fmt.Sprintf("token:%s", token)
	return s.client.Del(ctx, key).Err()
}

// AddToBlacklist 将 Token 加入黑名单
func (s *Store) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	key := fmt.Sprintf("token:blacklist:%s", token)
	return s.client.Set(ctx, key, "1", expiration).Err()
}

// IsInBlacklist 检查 Token 是否在黑名单中
func (s *Store) IsInBlacklist(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", token)
	count, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// SetRefreshToken 设置 Refresh Token
func (s *Store) SetRefreshToken(ctx context.Context, refreshToken, userID string, expiration time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.client.Set(ctx, key, userID, expiration).Err()
}

// GetRefreshToken 获取 Refresh Token 对应的用户ID
func (s *Store) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.client.Get(ctx, key).Result()
}

// DeleteRefreshToken 删除 Refresh Token
func (s *Store) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	key := fmt.Sprintf("refresh_token:%s", refreshToken)
	return s.client.Del(ctx, key).Err()
}
