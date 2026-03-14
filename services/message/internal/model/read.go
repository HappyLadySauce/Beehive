package model

import (
	"time"

	"gorm.io/gorm"
)

// ConversationRead 用户在某会话的已读位置（last_read_server_time）
// 表名 conversation_read，唯一约束 (user_id, conversation_id)；user_id 为 10 位，conversation_id 为 varchar(36)
type ConversationRead struct {
	UserID              string    `gorm:"column:user_id;type:char(10);primaryKey"`
	ConversationID      string    `gorm:"column:conversation_id;type:varchar(36);primaryKey"`
	LastReadServerTime  int64     `gorm:"column:last_read_server_time;not null;default:0"`
	UpdatedAt           time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (ConversationRead) TableName() string {
	return "conversation_read"
}

type ReadModel struct {
	db *gorm.DB
}

func NewReadModel(db *gorm.DB) *ReadModel {
	return &ReadModel{db: db}
}

// UpsertLastRead 更新或插入用户的已读位置（按 last_read_server_time）
func (m *ReadModel) UpsertLastRead(userID, conversationID string, lastReadServerTime int64) error {
	now := time.Now()
	var r ConversationRead
	err := m.db.Where("user_id = ? AND conversation_id = ?", userID, conversationID).First(&r).Error
	if err == gorm.ErrRecordNotFound {
		return m.db.Create(&ConversationRead{
			UserID:             userID,
			ConversationID:     conversationID,
			LastReadServerTime: lastReadServerTime,
			UpdatedAt:          now,
		}).Error
	}
	if err != nil {
		return err
	}
	return m.db.Model(&r).Updates(map[string]interface{}{
		"last_read_server_time": lastReadServerTime,
		"updated_at":            now,
	}).Error
}

// GetLastReadTimes 批量获取用户在若干会话的已读时间，返回 map[conversation_id]last_read_server_time，未记录则为 0
func (m *ReadModel) GetLastReadTimes(userID string, conversationIDs []string) (map[string]int64, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	var list []ConversationRead
	err := m.db.Where("user_id = ? AND conversation_id IN ?", userID, conversationIDs).Find(&list).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string]int64, len(conversationIDs))
	for _, cid := range conversationIDs {
		out[cid] = 0
	}
	for _, r := range list {
		out[r.ConversationID] = r.LastReadServerTime
	}
	return out, nil
}
