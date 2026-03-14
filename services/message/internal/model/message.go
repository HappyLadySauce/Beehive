package model

import (
	"time"

	"gorm.io/gorm"
)

// Message 对应 messages 表
type Message struct {
	ID              string    `gorm:"column:id;type:uuid;primaryKey"`
	ServerMsgID     string    `gorm:"column:server_msg_id;type:text;uniqueIndex;not null"`
	ClientMsgID     string    `gorm:"column:client_msg_id;type:text;not null"`
	ConversationID  string    `gorm:"column:conversation_id;type:uuid;not null"`
	FromUserID      string    `gorm:"column:from_user_id;type:uuid;not null"`
	ToUserID        string    `gorm:"column:to_user_id;type:uuid"`
	BodyType        string    `gorm:"column:body_type;type:text;not null"`
	BodyText        string    `gorm:"column:body_text;type:text;not null"`
	ServerTime      int64     `gorm:"column:server_time;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (Message) TableName() string {
	return "messages"
}

type MessageModel struct {
	db *gorm.DB
}

func NewMessageModel(db *gorm.DB) *MessageModel {
	return &MessageModel{db: db}
}

func (m *MessageModel) Create(msg *Message) error {
	return m.db.Create(msg).Error
}

func (m *MessageModel) GetHistory(conversationID string, beforeTime int64, limit int) ([]*Message, error) {
	q := m.db.Where("conversation_id = ?", conversationID).Order("server_time DESC")
	if beforeTime > 0 {
		q = q.Where("server_time < ?", beforeTime)
	}
	var list []*Message
	err := q.Limit(limit).Find(&list).Error
	return list, err
}

func (m *MessageModel) GetLastByConversations(conversationIDs []string) (map[string]*Message, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	// 每个会话取 server_time 最大的一条：用 DISTINCT ON (conversation_id) 或子查询
	var list []*Message
	err := m.db.Raw(`
		SELECT DISTINCT ON (conversation_id) id, server_msg_id, client_msg_id, conversation_id, from_user_id, to_user_id, body_type, body_text, server_time, created_at
		FROM messages
		WHERE conversation_id IN ?
		ORDER BY conversation_id, server_time DESC
	`, conversationIDs).Scan(&list).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]*Message, len(list))
	for _, msg := range list {
		out[msg.ConversationID] = msg
	}
	return out, nil
}
