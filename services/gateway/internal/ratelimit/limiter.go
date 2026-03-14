package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// MessageSendLimiter 按 userId 对 message.send 做固定窗口限流（每分钟 N 条）
type MessageSendLimiter struct {
	rdb   *redis.Client
	limit int
}

// NewMessageSendLimiter 创建限流器，limit 为每用户每分钟允许的条数
func NewMessageSendLimiter(rdb *redis.Client, limit int) *MessageSendLimiter {
	return &MessageSendLimiter{rdb: rdb, limit: limit}
}

// Allow 检查是否允许该用户发送一条消息；若超限返回 false
func (l *MessageSendLimiter) Allow(ctx context.Context, userID string) (bool, error) {
	if l.rdb == nil || l.limit <= 0 {
		return true, nil
	}
	window := time.Now().Unix() / 60
	key := fmt.Sprintf("rate:message.send:%s:%d", userID, window)
	pipe := l.rdb.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 90*time.Second) // 窗口 key 保留 90 秒
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return incr.Val() <= int64(l.limit), nil
}

// Close 关闭限流器持有的 Redis 连接池，应在进程退出或不再使用限流时调用
func (l *MessageSendLimiter) Close() error {
	if l.rdb == nil {
		return nil
	}
	return l.rdb.Close()
}
