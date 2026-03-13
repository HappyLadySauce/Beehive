// Package ws 提供 Gateway WebSocket 连接与消息协议，与 docs/API/websocket-client-api.md 对齐。
package ws

// Envelope 为 WebSocket 统一 JSON 消息外壳（请求/响应/推送）。
type Envelope struct {
	Type    string      `json:"type"`
	Tid     string      `json:"tid,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
	Error   *ErrBody    `json:"error,omitempty"`
}

// ErrBody 为响应中的错误信息。
type ErrBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
