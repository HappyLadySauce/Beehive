package websocket

import (
	"context"

	"k8s.io/klog/v2"

	messagepb "github.com/HappyLadySauce/Beehive/api/proto/message/v1"
	"github.com/HappyLadySauce/Beehive/internal/beehive-gateway/client"
)

// MessageHandler 消息处理器
type MessageHandler struct {
	grpcClient *client.Client
	pusher     *Pusher
}

// NewMessageHandler 创建新的消息处理器
func NewMessageHandler(grpcClient *client.Client, pusher *Pusher) *MessageHandler {
	return &MessageHandler{
		grpcClient: grpcClient,
		pusher:     pusher,
	}
}

// HandleMessage 处理客户端消息
func (h *MessageHandler) HandleMessage(ctx context.Context, userID string, msg Message) {
	switch msg.Type {
	case "send_message":
		h.handleSendMessage(ctx, userID, msg)
	case "send_group_message":
		h.handleSendGroupMessage(ctx, userID, msg)
	case "get_history":
		h.handleGetHistory(ctx, userID, msg)
	case "get_conversations":
		h.handleGetConversations(ctx, userID, msg)
	case "mark_as_read":
		h.handleMarkAsRead(ctx, userID, msg)
	default:
		h.pusher.SendError(userID, "UNKNOWN_MESSAGE_TYPE", "Unknown message type: "+msg.Type)
	}
}

// handleSendMessage 处理发送单聊消息
func (h *MessageHandler) handleSendMessage(ctx context.Context, userID string, msg Message) {
	toID, ok := msg.Data["to_id"].(string)
	if !ok || toID == "" {
		h.pusher.SendError(userID, "INVALID_REQUEST", "Missing or invalid to_id")
		return
	}

	content, ok := msg.Data["content"].(string)
	if !ok || content == "" {
		h.pusher.SendError(userID, "INVALID_REQUEST", "Missing or invalid content")
		return
	}

	messageType := "text"
	if mt, ok := msg.Data["message_type"].(string); ok && mt != "" {
		messageType = mt
	}

	req := &messagepb.SendMessageRequest{
		FromId:      userID,
		ToId:        toID,
		Content:     content,
		MessageType: messageType,
	}

	resp, err := h.grpcClient.MessageService().SendMessage(ctx, req)
	if err != nil {
		klog.Errorf("Failed to send message: %v", err)
		h.pusher.SendError(userID, "SEND_MESSAGE_FAILED", err.Error())
		return
	}

	h.pusher.SendSuccess(userID, "message_sent", map[string]interface{}{
		"message_id": resp.MessageId,
		"created_at": resp.CreatedAt,
	})
}

// handleSendGroupMessage 处理发送群聊消息
func (h *MessageHandler) handleSendGroupMessage(ctx context.Context, userID string, msg Message) {
	groupID, ok := msg.Data["group_id"].(string)
	if !ok || groupID == "" {
		h.pusher.SendError(userID, "INVALID_REQUEST", "Missing or invalid group_id")
		return
	}

	content, ok := msg.Data["content"].(string)
	if !ok || content == "" {
		h.pusher.SendError(userID, "INVALID_REQUEST", "Missing or invalid content")
		return
	}

	messageType := "text"
	if mt, ok := msg.Data["message_type"].(string); ok && mt != "" {
		messageType = mt
	}

	req := &messagepb.SendGroupMessageRequest{
		FromId:      userID,
		GroupId:     groupID,
		Content:     content,
		MessageType: messageType,
	}

	resp, err := h.grpcClient.MessageService().SendGroupMessage(ctx, req)
	if err != nil {
		klog.Errorf("Failed to send group message: %v", err)
		h.pusher.SendError(userID, "SEND_GROUP_MESSAGE_FAILED", err.Error())
		return
	}

	h.pusher.SendSuccess(userID, "group_message_sent", map[string]interface{}{
		"message_id": resp.MessageId,
		"created_at": resp.CreatedAt,
	})
}

// handleGetHistory 处理获取消息历史
func (h *MessageHandler) handleGetHistory(ctx context.Context, userID string, msg Message) {
	req := &messagepb.GetMessageHistoryRequest{
		Id: userID,
	}

	if targetUserID, ok := msg.Data["target_user_id"].(string); ok && targetUserID != "" {
		req.TargetUserId = targetUserID
	}

	if groupID, ok := msg.Data["group_id"].(string); ok && groupID != "" {
		req.GroupId = groupID
	}

	limit := int32(50)
	if l, ok := msg.Data["limit"].(float64); ok {
		limit = int32(l)
	}
	req.Limit = limit

	offset := int32(0)
	if o, ok := msg.Data["offset"].(float64); ok {
		offset = int32(o)
	}
	req.Offset = offset

	resp, err := h.grpcClient.MessageService().GetMessageHistory(ctx, req)
	if err != nil {
		klog.Errorf("Failed to get message history: %v", err)
		h.pusher.SendError(userID, "GET_HISTORY_FAILED", err.Error())
		return
	}

	// 转换消息格式
	messages := make([]map[string]interface{}, len(resp.Messages))
	for i, m := range resp.Messages {
		messages[i] = map[string]interface{}{
			"id":           m.Id,
			"type":         m.Type,
			"from_id":      m.FromId,
			"to_id":        m.ToId,
			"group_id":     m.GroupId,
			"content":      m.Content,
			"message_type": m.MessageType,
			"status":       m.Status,
			"created_at":   m.CreatedAt,
			"updated_at":   m.UpdatedAt,
		}
	}

	h.pusher.SendSuccess(userID, "message_history", map[string]interface{}{
		"messages": messages,
		"total":    resp.Total,
	})
}

// handleGetConversations 处理获取会话列表
func (h *MessageHandler) handleGetConversations(ctx context.Context, userID string, msg Message) {
	req := &messagepb.GetConversationsRequest{
		Id: userID,
	}

	limit := int32(50)
	if l, ok := msg.Data["limit"].(float64); ok {
		limit = int32(l)
	}
	req.Limit = limit

	offset := int32(0)
	if o, ok := msg.Data["offset"].(float64); ok {
		offset = int32(o)
	}
	req.Offset = offset

	resp, err := h.grpcClient.MessageService().GetConversations(ctx, req)
	if err != nil {
		klog.Errorf("Failed to get conversations: %v", err)
		h.pusher.SendError(userID, "GET_CONVERSATIONS_FAILED", err.Error())
		return
	}

	// 转换会话格式
	conversations := make([]map[string]interface{}, len(resp.Conversations))
	for i, c := range resp.Conversations {
		conversations[i] = map[string]interface{}{
			"id":              c.Id,
			"type":            c.Type,
			"id1":             c.Id1,
			"id2":             c.Id2,
			"group_id":        c.GroupId,
			"last_message_id": c.LastMessageId,
			"last_message_at": c.LastMessageAt,
			"unread_count":    c.UnreadCount,
			"created_at":      c.CreatedAt,
			"updated_at":      c.UpdatedAt,
		}
	}

	h.pusher.SendSuccess(userID, "conversations", map[string]interface{}{
		"conversations": conversations,
		"total":         resp.Total,
	})
}

// handleMarkAsRead 处理标记消息已读
func (h *MessageHandler) handleMarkAsRead(ctx context.Context, userID string, msg Message) {
	req := &messagepb.MarkAsReadRequest{
		Id: userID,
	}

	if messageIDs, ok := msg.Data["message_ids"].([]interface{}); ok {
		ids := make([]string, len(messageIDs))
		for i, id := range messageIDs {
			if idStr, ok := id.(string); ok {
				ids[i] = idStr
			}
		}
		req.MessageIds = ids
	}

	if conversationID, ok := msg.Data["conversation_id"].(string); ok && conversationID != "" {
		req.ConversationId = conversationID
	}

	resp, err := h.grpcClient.MessageService().MarkAsRead(ctx, req)
	if err != nil {
		klog.Errorf("Failed to mark messages as read: %v", err)
		h.pusher.SendError(userID, "MARK_AS_READ_FAILED", err.Error())
		return
	}

	h.pusher.SendSuccess(userID, "messages_marked_read", map[string]interface{}{
		"success": resp.Success,
	})
}
