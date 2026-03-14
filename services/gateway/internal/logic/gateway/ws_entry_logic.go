// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ws"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type WsEntryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWsEntryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WsEntryLogic {
	return &WsEntryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ServeConn 在连接上运行读循环：解析 JSON Envelope，按 type 路由（为 auth/presence 等集成做准备）。
// 连接关闭由调用方（Handler）负责。
func (l *WsEntryLogic) ServeConn(c *ws.Connection) error {
	conn := c.Conn()
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				return nil
			}
			return err
		}
		var env ws.Envelope
		if err := json.Unmarshal(data, &env); err != nil {
			l.sendError(c, env.Tid, "bad_request", "invalid json envelope")
			continue
		}
		if env.Type == "" {
			l.sendError(c, env.Tid, "bad_request", "missing type")
			continue
		}
		l.dispatch(c, &env)
	}
}

func (l *WsEntryLogic) dispatch(c *ws.Connection, env *ws.Envelope) {
	// 除 auth.* 之外的消息都要求连接已登录，包括 presence.ping 在内。
	if env.Type != "auth.login" && env.Type != "auth.tokenLogin" && env.Type != "auth.logout" {
		if c.UserID == "" {
			l.sendError(c, env.Tid, "unauthorized", "user not logged in")
			return
		}
	}

	switch env.Type {
	case "presence.ping":
		l.handlePresencePing(c, env)
	case "auth.login", "auth.tokenLogin":
		l.handleAuth(c, env)
	case "auth.logout":
		l.handleAuthLogout(c, env)
	default:
		l.sendError(c, env.Tid, "bad_request", "unknown type: "+env.Type)
	}
}

func (l *WsEntryLogic) handlePresencePing(c *ws.Connection, env *ws.Envelope) {
	if c.UserID != "" {
		_, err := l.svcCtx.PresenceSvc.RefreshSession(l.ctx, &presenceservice.RefreshSessionRequest{
			UserId: c.UserID,
			ConnId: c.ConnID,
		})
		if err != nil {
			l.Errorf("refresh session failed for user %s conn %s: %v", c.UserID, c.ConnID, err)
		}
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "presence.ping.ok",
		Tid:     env.Tid,
		Payload: map[string]int64{"serverTime": time.Now().Unix()},
		Error:   nil,
	})
}

func (l *WsEntryLogic) sendError(c *ws.Connection, tid, code, message string) {
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "",
		Tid:     tid,
		Payload: nil,
		Error:   &ws.ErrBody{Code: code, Message: message},
	})
}

// handleAuth 处理 auth.login / auth.tokenLogin。
func (l *WsEntryLogic) handleAuth(c *ws.Connection, env *ws.Envelope) {
	switch env.Type {
	case "auth.login":
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
			DeviceID string `json:"deviceId"`
		}
		if !l.bindJSONPayload(c, env, &payload) {
			return
		}
		resp, err := l.svcCtx.AuthSvc.Login(l.ctx, &authservice.LoginRequest{
			Username: payload.Username,
			Password: payload.Password,
			DeviceId: payload.DeviceID,
		})
		if err != nil {
			l.sendError(c, env.Tid, "unauthorized", err.Error())
			return
		}
		l.afterAuthSuccess(c, env, resp, payload.DeviceID)
	case "auth.tokenLogin":
		var payload struct {
			AccessToken string `json:"accessToken"`
			DeviceID    string `json:"deviceId"`
		}
		if !l.bindJSONPayload(c, env, &payload) {
			return
		}
		resp, err := l.svcCtx.AuthSvc.TokenLogin(l.ctx, &authservice.TokenLoginRequest{
			AccessToken: payload.AccessToken,
			DeviceId:    payload.DeviceID,
		})
		if err != nil {
			l.sendError(c, env.Tid, "unauthorized", err.Error())
			return
		}
		l.afterAuthSuccess(c, env, resp, payload.DeviceID)
	default:
		l.sendError(c, env.Tid, "bad_request", "unsupported auth type")
	}
}

// handleAuthLogout 处理 auth.logout。
func (l *WsEntryLogic) handleAuthLogout(c *ws.Connection, env *ws.Envelope) {
	var payload struct {
		AccessToken string `json:"accessToken"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.AccessToken != "" {
		_, _ = l.svcCtx.AuthSvc.Logout(l.ctx, &authservice.LogoutRequest{
			AccessToken: payload.AccessToken,
		})
	}
	if c.UserID != "" {
		_, _ = l.svcCtx.PresenceSvc.UnregisterSession(l.ctx, &presenceservice.UnregisterSessionRequest{
			UserId: c.UserID,
			ConnId: c.ConnID,
		})
		c.BindUser("")
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "auth.logout.ok",
		Tid:     env.Tid,
		Payload: map[string]any{"ok": true},
		Error:   nil,
	})
}

// afterAuthSuccess 在登录或 tokenLogin 成功后绑定 UserID 并注册 Presence。
func (l *WsEntryLogic) afterAuthSuccess(c *ws.Connection, env *ws.Envelope, resp *authservice.LoginResponse, deviceID string) {
	if resp == nil || resp.UserId == "" {
		l.sendError(c, env.Tid, "internal_error", "empty auth response")
		return
	}
	// 先向 Presence 注册会话，确保在线状态成功记录；失败则不认为登录成功。
	if _, err := l.svcCtx.PresenceSvc.RegisterSession(l.ctx, &presenceservice.RegisterSessionRequest{
		UserId:     resp.UserId,
		GatewayId:  l.svcCtx.Config.GatewayID,
		ConnId:     c.ConnID,
		DeviceId:   deviceID,
		DeviceType: "", // 可根据实际需求从 Query 或 payload 补充
		Ip:         "",
	}); err != nil {
		l.Errorf("register session failed for user %s conn %s: %v", resp.UserId, c.ConnID, err)
		l.sendError(c, env.Tid, "internal_error", "failed to register session")
		return
	}
	// Presence 注册成功后再在本地绑定用户，并返回登录成功响应。
	c.BindUser(resp.UserId)
	// 返回 auth.login.ok 或 auth.tokenLogin.ok
	okType := env.Type + ".ok"
	_ = c.WriteJSON(&ws.Envelope{
		Type: okType,
		Tid:  env.Tid,
		Payload: map[string]any{
			"userId":       resp.UserId,
			"accessToken":  resp.AccessToken,
			"refreshToken": resp.RefreshToken,
			"expiresIn":    resp.ExpiresIn,
		},
		Error: nil,
	})
}

// bindJSONPayload 将 Envelope.Payload 反序列化为给定结构体。
func (l *WsEntryLogic) bindJSONPayload(c *ws.Connection, env *ws.Envelope, v interface{}) bool {
	if env.Payload == nil {
		l.sendError(c, env.Tid, "bad_request", "missing payload")
		return false
	}
	raw, err := json.Marshal(env.Payload)
	if err != nil {
		l.sendError(c, env.Tid, "bad_request", "invalid payload format")
		return false
	}
	if err := json.Unmarshal(raw, v); err != nil {
		l.sendError(c, env.Tid, "bad_request", "invalid payload json")
		return false
	}
	return true
}
