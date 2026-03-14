package admin

import "time"

// formatUnixTime 将 Unix 时间戳（秒）格式化为 ISO8601 字符串；若 <= 0 返回空字符串。
func formatUnixTime(sec int64) string {
	if sec <= 0 {
		return ""
	}
	return time.Unix(sec, 0).UTC().Format(time.RFC3339)
}
