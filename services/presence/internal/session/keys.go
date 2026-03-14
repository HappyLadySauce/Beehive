package session

// Redis key 与 TTL 约定（与 plan 一致）

const (
	userConnsKeyPrefix = "presence:user:"
	userConnsKeySuffix = ":conns"
	sessionKeyPrefix   = "presence:session:"
)

// UserConnsKey 返回用户会话 connId 集合的 key：presence:user:{userId}:conns
func UserConnsKey(userID string) string {
	return userConnsKeyPrefix + userID + userConnsKeySuffix
}

// SessionKey 返回单条会话 Hash 的 key：presence:session:{userId}:{connId}
func SessionKey(userID, connID string) string {
	return sessionKeyPrefix + userID + ":" + connID
}

// Hash 字段名（与 logic 中 HSet 一致）
const (
	HashGatewayID   = "gateway_id"
	HashConnID      = "conn_id"
	HashDeviceID    = "device_id"
	HashDeviceType  = "device_type"
	HashLastPingAt  = "last_ping_at"
)

// DefaultSessionTTLSeconds 配置未设置或 <=0 时使用的默认会话 TTL（秒）
const DefaultSessionTTLSeconds = 90

// SessionTTL 返回会话 TTL 秒数；configSeconds <= 0 时使用默认值。
func SessionTTL(configSeconds int) int {
	if configSeconds <= 0 {
		return DefaultSessionTTLSeconds
	}
	return configSeconds
}
