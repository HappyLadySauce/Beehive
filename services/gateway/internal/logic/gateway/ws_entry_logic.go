// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package gateway

import (
	"context"
	"encoding/json"
	"time"

	"github.com/HappyLadySauce/Beehive/services/auth/authservice"
	"github.com/HappyLadySauce/Beehive/services/conversation/conversationservice"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/svc"
	"github.com/HappyLadySauce/Beehive/services/gateway/internal/ws"
	"github.com/HappyLadySauce/Beehive/services/message/messageservice"
	"github.com/HappyLadySauce/Beehive/services/presence/presenceservice"
	"github.com/HappyLadySauce/Beehive/services/user/userservice"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if env.Type != "auth.login" && env.Type != "auth.tokenLogin" && env.Type != "auth.logout" && env.Type != "auth.register" {
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
	case "auth.register":
		l.handleAuthRegister(c, env)
	case "auth.logout":
		l.handleAuthLogout(c, env)
	case "user.me":
		l.handleUserMe(c, env)
	case "conversation.list":
		l.handleConversationList(c, env)
	case "conversation.create":
		l.handleConversationCreate(c, env)
	case "conversation.addMember":
		l.handleConversationAddMember(c, env)
	case "conversation.removeMember":
		l.handleConversationRemoveMember(c, env)
	case "message.send":
		l.handleMessageSend(c, env)
	case "message.history":
		l.handleMessageHistory(c, env)
	case "message.read":
		l.handleMessageRead(c, env)
	case "contact.list":
		l.handleContactList(c, env)
	case "contact.add":
		l.handleContactAdd(c, env)
	case "contact.remove":
		l.handleContactRemove(c, env)
	default:
		l.sendError(c, env.Tid, "bad_request", "unknown type: "+env.Type)
	}
}

func (l *WsEntryLogic) handlePresencePing(c *ws.Connection, env *ws.Envelope) {
	if c.UserID != "" && l.svcCtx.PresenceSvc != nil {
		l.Infow("rpc call", logx.Field("method", "presence.RefreshSession"), logx.Field("userId", c.UserID), logx.Field("connId", c.ConnID))
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

// handleUserMe 获取当前登录用户资料，需已登录；调用 UserService.GetUser。未配置 UserService 时返回不可用。
func (l *WsEntryLogic) handleUserMe(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.UserSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "user service not configured")
		return
	}
	l.Infow("rpc call", logx.Field("method", "user.GetUser"), logx.Field("userId", c.UserID))
	resp, err := l.svcCtx.UserSvc.GetUser(l.ctx, &userservice.GetUserRequest{Id: c.UserID})
	if err != nil {
		l.Errorf("get user me failed for %s: %v", c.UserID, err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	if resp.GetUser() == nil {
		l.sendError(c, env.Tid, "not_found", "user not found")
		return
	}
	u := resp.User
	_ = c.WriteJSON(&ws.Envelope{
		Type: env.Type + ".ok",
		Tid:  env.Tid,
		Payload: map[string]any{
			"id":        u.GetId(),
			"nickname":  u.GetNickname(),
			"avatarUrl": u.GetAvatarUrl(),
			"bio":       u.GetBio(),
			"status":    u.GetStatus(),
		},
		Error: nil,
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

// handleAuthRegister 处理 auth.register，调用 AuthService.Register，成功后与登录一致绑定并返回 token。
func (l *WsEntryLogic) handleAuthRegister(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.AuthSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "auth service not configured")
		return
	}
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	l.Infow("rpc call", logx.Field("method", "auth.Register"), logx.Field("username", payload.Username))
	resp, err := l.svcCtx.AuthSvc.Register(l.ctx, &authservice.RegisterRequest{
		Username: payload.Username,
		Password: payload.Password,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.AlreadyExists:
				l.sendError(c, env.Tid, "bad_request", st.Message())
				return
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", st.Message())
				return
			}
		}
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	l.afterAuthSuccess(c, env, resp, "")
}

// handleAuth 处理 auth.login / auth.tokenLogin。
func (l *WsEntryLogic) handleAuth(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.AuthSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "auth service not configured")
		return
	}
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
		l.Infow("rpc call", logx.Field("method", "auth.Login"), logx.Field("username", payload.Username))
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
		l.Infow("rpc call", logx.Field("method", "auth.TokenLogin"))
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
	if l.svcCtx.AuthSvc != nil && payload.AccessToken != "" {
		l.Infow("rpc call", logx.Field("method", "auth.Logout"))
		_, _ = l.svcCtx.AuthSvc.Logout(l.ctx, &authservice.LogoutRequest{
			AccessToken: payload.AccessToken,
		})
	}
	if l.svcCtx.PresenceSvc != nil && c.UserID != "" {
		l.Infow("rpc call", logx.Field("method", "presence.UnregisterSession"), logx.Field("userId", c.UserID), logx.Field("connId", c.ConnID))
		_, _ = l.svcCtx.PresenceSvc.UnregisterSession(l.ctx, &presenceservice.UnregisterSessionRequest{
			UserId: c.UserID,
			ConnId: c.ConnID,
		})
	}
	if c.UserID != "" {
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
	if l.svcCtx.PresenceSvc != nil {
		l.Infow("rpc call", logx.Field("method", "presence.RegisterSession"), logx.Field("userId", resp.UserId), logx.Field("connId", c.ConnID))
		if _, err := l.svcCtx.PresenceSvc.RegisterSession(l.ctx, &presenceservice.RegisterSessionRequest{
			UserId:     resp.UserId,
			GatewayId:  l.svcCtx.Config.GatewayID,
			ConnId:     c.ConnID,
			DeviceId:   deviceID,
			DeviceType: "",
			Ip:         "",
		}); err != nil {
			l.Errorf("register session failed for user %s conn %s: %v", resp.UserId, c.ConnID, err)
			l.sendError(c, env.Tid, "internal_error", "failed to register session")
			return
		}
	}
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

func (l *WsEntryLogic) handleConversationList(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.ConversationSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "conversation service not configured")
		return
	}
	var payload struct {
		Cursor string `json:"cursor"`
		Limit  int32  `json:"limit"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.Limit <= 0 {
		payload.Limit = 50
	}
	if payload.Limit > 100 {
		payload.Limit = 100
	}
	resp, err := l.svcCtx.ConversationSvc.ListUserConversations(l.ctx, &conversationservice.ListUserConversationsRequest{
		UserId: c.UserID,
		Cursor: payload.Cursor,
		Limit:  payload.Limit,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.InvalidArgument {
			l.sendError(c, env.Tid, "bad_request", s.Message())
			return
		}
		l.Errorf("list user conversations failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	convIDs := make([]string, 0, len(resp.Items))
	for _, item := range resp.Items {
		convIDs = append(convIDs, item.Id)
	}
	var lastMessages map[string]*messageservice.MessageRecord
	var unreadCounts map[string]int32
	if l.svcCtx.MessageSvc != nil && len(convIDs) > 0 {
		lastResp, _ := l.svcCtx.MessageSvc.GetLastMessages(l.ctx, &messageservice.GetLastMessagesRequest{ConversationIds: convIDs})
		if lastResp != nil {
			lastMessages = lastResp.LastMessages
		}
		unreadResp, _ := l.svcCtx.MessageSvc.GetUnreadCounts(l.ctx, &messageservice.GetUnreadCountsRequest{UserId: c.UserID, ConversationIds: convIDs})
		if unreadResp != nil && unreadResp.Counts != nil {
			unreadCounts = unreadResp.Counts
		}
	}
	items := make([]map[string]any, 0, len(resp.Items))
	for _, item := range resp.Items {
		unread := int(0)
		if unreadCounts != nil {
			if n, ok := unreadCounts[item.Id]; ok {
				unread = int(n)
			}
		}
		entry := map[string]any{
			"id":            item.Id,
			"name":          item.Name,
			"avatar":        "",
			"type":          item.Type,
			"unreadCount":   unread,
			"lastActiveAt":  item.LastActiveAt,
		}
		if lastMessages != nil {
			if lm, ok := lastMessages[item.Id]; ok && lm != nil {
				preview := ""
				if lm.Body != nil {
					preview = lm.Body.Text
				}
				entry["lastMessage"] = map[string]any{
					"serverMsgId": lm.ServerMsgId,
					"preview":     preview,
					"serverTime":  lm.ServerTime,
				}
			}
		}
		items = append(items, entry)
	}
	var nextCursor interface{}
	if resp.NextCursor != "" {
		nextCursor = resp.NextCursor
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "conversation.list.ok",
		Tid:     env.Tid,
		Payload: map[string]any{"items": items, "nextCursor": nextCursor},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleConversationCreate(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.ConversationSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "conversation service not configured")
		return
	}
	var payload struct {
		Type       string   `json:"type"`
		Name       string   `json:"name"`
		MemberIds  []string `json:"memberIds"`
		ToUsername string   `json:"toUsername"`
		ToAccount  string   `json:"toAccount"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	memberIds := payload.MemberIds
	if payload.Type == "single" && len(memberIds) == 0 && (payload.ToUsername != "" || payload.ToAccount != "") {
		if l.svcCtx.UserSvc == nil {
			l.sendError(c, env.Tid, "unavailable", "user service not configured for resolve")
			return
		}
		var toUserId string
		if payload.ToAccount != "" {
			u, err := l.svcCtx.UserSvc.GetUser(l.ctx, &userservice.GetUserRequest{Id: payload.ToAccount})
			if err != nil || u.GetUser() == nil {
				l.sendError(c, env.Tid, "not_found", "user not found")
				return
			}
			toUserId = u.User.Id
		} else {
			u, err := l.svcCtx.UserSvc.GetUserByUsername(l.ctx, &userservice.GetUserByUsernameRequest{Username: payload.ToUsername})
			if err != nil || u.GetUser() == nil {
				l.sendError(c, env.Tid, "not_found", "user not found")
				return
			}
			toUserId = u.User.Id
		}
		if toUserId == c.UserID {
			l.sendError(c, env.Tid, "bad_request", "cannot create chat with self")
			return
		}
		memberIds = []string{c.UserID, toUserId}
	}
	resp, err := l.svcCtx.ConversationSvc.CreateConversation(l.ctx, &conversationservice.CreateConversationRequest{
		Type:      payload.Type,
		Name:      payload.Name,
		MemberIds: memberIds,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.AlreadyExists:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			}
		}
		l.Errorf("create conversation failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "conversation.create.ok",
		Tid:     env.Tid,
		Payload: map[string]any{"conversationId": resp.ConversationId},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleConversationAddMember(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.ConversationSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "conversation service not configured")
		return
	}
	var payload struct {
		ConversationId string `json:"conversationId"`
		UserId         string `json:"userId"`
		Role           string `json:"role"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.ConversationId == "" || payload.UserId == "" {
		l.sendError(c, env.Tid, "bad_request", "conversationId and userId are required")
		return
	}
	if payload.Role == "" {
		payload.Role = "member"
	}
	_, err := l.svcCtx.ConversationSvc.AddMember(l.ctx, &conversationservice.AddMemberRequest{
		ConversationId: payload.ConversationId,
		UserId:         payload.UserId,
		Role:           payload.Role,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.NotFound:
				l.sendError(c, env.Tid, "not_found", s.Message())
				return
			case codes.AlreadyExists:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			}
		}
		l.Errorf("add member failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "conversation.addMember.ok",
		Tid:     env.Tid,
		Payload: map[string]any{},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleConversationRemoveMember(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.ConversationSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "conversation service not configured")
		return
	}
	var payload struct {
		ConversationId string `json:"conversationId"`
		UserId         string `json:"userId"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.ConversationId == "" || payload.UserId == "" {
		l.sendError(c, env.Tid, "bad_request", "conversationId and userId are required")
		return
	}
	_, err := l.svcCtx.ConversationSvc.RemoveMember(l.ctx, &conversationservice.RemoveMemberRequest{
		ConversationId: payload.ConversationId,
		UserId:         payload.UserId,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.NotFound:
				l.sendError(c, env.Tid, "not_found", s.Message())
				return
			}
		}
		l.Errorf("remove member failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "conversation.removeMember.ok",
		Tid:     env.Tid,
		Payload: map[string]any{},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleMessageSend(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.MessageSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "message service not configured")
		return
	}
	if l.svcCtx.MessageSendLimit != nil {
		allow, err := l.svcCtx.MessageSendLimit.Allow(l.ctx, c.UserID)
		if err != nil {
			l.Errorf("rate limit check failed: %v", err)
			l.sendError(c, env.Tid, "internal_error", "rate limit check failed")
			return
		}
		if !allow {
			_ = c.WriteJSON(&ws.Envelope{
				Type:    "message.send.error",
				Tid:     env.Tid,
				Payload: nil,
				Error:   &ws.ErrBody{Code: "rate_limited", Message: "too many requests"},
			})
			return
		}
	}
	var payload struct {
		ClientMsgId     string `json:"clientMsgId"`
		ConversationId string `json:"conversationId"`
		ToUserId        string `json:"toUserId"`
		Body            struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"body"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	conversationId := payload.ConversationId
	if conversationId == "" {
		if payload.ToUserId == "" {
			l.sendError(c, env.Tid, "bad_request", "conversationId or toUserId is required")
			return
		}
		if l.svcCtx.ConversationSvc == nil {
			l.sendError(c, env.Tid, "unavailable", "conversation service not configured for single chat")
			return
		}
		findResp, err := l.svcCtx.ConversationSvc.FindOrCreateSingleConversation(l.ctx, &conversationservice.FindOrCreateSingleConversationRequest{
			UserId_1: c.UserID,
			UserId_2: payload.ToUserId,
		})
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.InvalidArgument {
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			}
			l.Errorf("find or create single conversation failed: %v", err)
			l.sendError(c, env.Tid, "internal_error", err.Error())
			return
		}
		conversationId = findResp.ConversationId
	}
	if payload.Body.Type == "" {
		payload.Body.Type = "text"
	}
	resp, err := l.svcCtx.MessageSvc.PostMessage(l.ctx, &messageservice.PostMessageRequest{
		ClientMsgId:     payload.ClientMsgId,
		ConversationId:  conversationId,
		FromUserId:     c.UserID,
		ToUserId:       payload.ToUserId,
		Body:           &messageservice.MessageBody{Type: payload.Body.Type, Text: payload.Body.Text},
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.NotFound:
				l.sendError(c, env.Tid, "not_found", s.Message())
				return
			}
		}
		l.Errorf("post message failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "message.send.ok",
		Tid:     env.Tid,
		Payload: map[string]any{
			"serverMsgId":    resp.ServerMsgId,
			"serverTime":     resp.ServerTime,
			"conversationId": resp.ConversationId,
		},
		Error: nil,
	})
}

func (l *WsEntryLogic) handleMessageHistory(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.MessageSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "message service not configured")
		return
	}
	var payload struct {
		ConversationId string `json:"conversationId"`
		Before         int64  `json:"before"`
		Limit          int32  `json:"limit"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.ConversationId == "" {
		l.sendError(c, env.Tid, "bad_request", "conversationId is required")
		return
	}
	if payload.Limit <= 0 {
		payload.Limit = 50
	}
	if payload.Limit > 100 {
		payload.Limit = 100
	}
	resp, err := l.svcCtx.MessageSvc.GetHistory(l.ctx, &messageservice.GetHistoryRequest{
		ConversationId: payload.ConversationId,
		BeforeTime:    payload.Before,
		Limit:         payload.Limit,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok && s.Code() == codes.InvalidArgument {
			l.sendError(c, env.Tid, "bad_request", s.Message())
			return
		}
		l.Errorf("get history failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	items := make([]map[string]any, 0, len(resp.Items))
	for _, m := range resp.Items {
		body := map[string]any{"type": "text"}
		if m.Body != nil {
			body["type"] = m.Body.Type
			body["text"] = m.Body.Text
		}
		items = append(items, map[string]any{
			"serverMsgId":    m.ServerMsgId,
			"clientMsgId":    nil,
			"conversationId": m.ConversationId,
			"fromUserId":     m.FromUserId,
			"toUserId":       m.ToUserId,
			"body":           body,
			"serverTime":     m.ServerTime,
		})
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "message.history.ok",
		Tid:     env.Tid,
		Payload: map[string]any{"items": items, "hasMore": resp.HasMore},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleMessageRead(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.MessageSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "message service not configured")
		return
	}
	var payload struct {
		ConversationId string `json:"conversationId"`
		ServerMsgId    string `json:"serverMsgId"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.ConversationId == "" || payload.ServerMsgId == "" {
		l.sendError(c, env.Tid, "bad_request", "conversationId and serverMsgId are required")
		return
	}
	_, err := l.svcCtx.MessageSvc.MarkRead(l.ctx, &messageservice.MarkReadRequest{
		UserId:         c.UserID,
		ConversationId:  payload.ConversationId,
		ServerMsgId:     payload.ServerMsgId,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.NotFound:
				l.sendError(c, env.Tid, "not_found", s.Message())
				return
			}
		}
		l.Errorf("mark read failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "message.read.ok",
		Tid:     env.Tid,
		Payload: map[string]any{},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleContactList(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.UserSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "user service not configured")
		return
	}
	resp, err := l.svcCtx.UserSvc.ListContacts(l.ctx, &userservice.ListContactsRequest{OwnerId: c.UserID})
	if err != nil {
		l.Errorf("list contacts failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "contact.list.ok",
		Tid:     env.Tid,
		Payload: map[string]any{"contactUserIds": resp.ContactUserIds},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleContactAdd(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.UserSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "user service not configured")
		return
	}
	var payload struct {
		ToUserId   string `json:"toUserId"`
		ToUsername string `json:"toUsername"`
		ToAccount  string `json:"toAccount"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	contactUserID := payload.ToUserId
	if contactUserID == "" && payload.ToAccount != "" {
		u, err := l.svcCtx.UserSvc.GetUser(l.ctx, &userservice.GetUserRequest{Id: payload.ToAccount})
		if err != nil || u.GetUser() == nil {
			l.sendError(c, env.Tid, "not_found", "user not found")
			return
		}
		contactUserID = u.User.Id
	}
	if contactUserID == "" && payload.ToUsername != "" {
		u, err := l.svcCtx.UserSvc.GetUserByUsername(l.ctx, &userservice.GetUserByUsernameRequest{Username: payload.ToUsername})
		if err != nil || u.GetUser() == nil {
			l.sendError(c, env.Tid, "not_found", "user not found")
			return
		}
		contactUserID = u.User.Id
	}
	if contactUserID == "" {
		l.sendError(c, env.Tid, "bad_request", "toUserId, toUsername or toAccount is required")
		return
	}
	_, err := l.svcCtx.UserSvc.AddContact(l.ctx, &userservice.AddContactRequest{
		OwnerId:       c.UserID,
		ContactUserId: contactUserID,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			switch s.Code() {
			case codes.InvalidArgument:
				l.sendError(c, env.Tid, "bad_request", s.Message())
				return
			case codes.NotFound:
				l.sendError(c, env.Tid, "not_found", s.Message())
				return
			}
		}
		l.Errorf("add contact failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "contact.add.ok",
		Tid:     env.Tid,
		Payload: map[string]any{},
		Error:   nil,
	})
}

func (l *WsEntryLogic) handleContactRemove(c *ws.Connection, env *ws.Envelope) {
	if l.svcCtx.UserSvc == nil {
		l.sendError(c, env.Tid, "unavailable", "user service not configured")
		return
	}
	var payload struct {
		ContactUserId string `json:"contactUserId"`
	}
	if !l.bindJSONPayload(c, env, &payload) {
		return
	}
	if payload.ContactUserId == "" {
		l.sendError(c, env.Tid, "bad_request", "contactUserId is required")
		return
	}
	_, err := l.svcCtx.UserSvc.RemoveContact(l.ctx, &userservice.RemoveContactRequest{
		OwnerId:       c.UserID,
		ContactUserId: payload.ContactUserId,
	})
	if err != nil {
		l.Errorf("remove contact failed: %v", err)
		l.sendError(c, env.Tid, "internal_error", err.Error())
		return
	}
	_ = c.WriteJSON(&ws.Envelope{
		Type:    "contact.remove.ok",
		Tid:     env.Tid,
		Payload: map[string]any{},
		Error:   nil,
	})
}
