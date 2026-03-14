package logic

import (
	"context"
	"strconv"

	"github.com/HappyLadySauce/Beehive/services/presence/internal/session"
	"github.com/HappyLadySauce/Beehive/services/presence/pb"
	"github.com/redis/go-redis/v9"
)

// getSessionsForUser 从 Redis 读取某用户的在线会话列表，并清理已过期的 connId 索引。
func getSessionsForUser(ctx context.Context, rdb *redis.Client, userID string) ([]*pb.SessionInfo, error) {
	if userID == "" {
		return nil, nil
	}
	userConnsKey := session.UserConnsKey(userID)
	connIds, err := rdb.SMembers(ctx, userConnsKey).Result()
	if err != nil {
		return nil, err
	}
	var sessions []*pb.SessionInfo
	for _, connID := range connIds {
		sessionKey := session.SessionKey(userID, connID)
		m, err := rdb.HGetAll(ctx, sessionKey).Result()
		if err != nil {
			continue
		}
		if len(m) == 0 {
			rdb.SRem(ctx, userConnsKey, connID)
			continue
		}
		var lastPingAt int64
		if s := m[session.HashLastPingAt]; s != "" {
			lastPingAt, _ = strconv.ParseInt(s, 10, 64)
		}
		sessions = append(sessions, &pb.SessionInfo{
			GatewayId:   m[session.HashGatewayID],
			ConnId:      m[session.HashConnID],
			DeviceId:    m[session.HashDeviceID],
			DeviceType:  m[session.HashDeviceType],
			LastPingAt:  lastPingAt,
		})
	}
	return sessions, nil
}
