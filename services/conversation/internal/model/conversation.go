package model

import (
	"time"

	"gorm.io/gorm"
)

// Conversation 对应 conversations 表
type Conversation struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	Type         string    `gorm:"column:type;type:text;not null;default:single"`
	Name         string    `gorm:"column:name;type:text;not null;default:''"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	LastActiveAt time.Time `gorm:"column:last_active_at;type:timestamptz;not null"`
}

func (Conversation) TableName() string {
	return "conversations"
}

// ConversationMember 对应 conversation_members 表
type ConversationMember struct {
	ID             string    `gorm:"column:id;type:uuid;primaryKey"`
	ConversationID string    `gorm:"column:conversation_id;type:uuid;not null;uniqueIndex:uq_conv_member"`
	UserID         string    `gorm:"column:user_id;type:uuid;not null;uniqueIndex:uq_conv_member"`
	Role           string    `gorm:"column:role;type:text;not null;default:member"`
	Status         string    `gorm:"column:status;type:text;not null;default:active"`
	JoinedAt       time.Time `gorm:"column:joined_at;type:timestamptz;not null"`
}

func (ConversationMember) TableName() string {
	return "conversation_members"
}

type ConversationModel struct {
	db *gorm.DB
}

func NewConversationModel(db *gorm.DB) *ConversationModel {
	return &ConversationModel{db: db}
}

// Create 创建会话并批量插入成员（同一事务）
func (m *ConversationModel) Create(c *Conversation, members []*ConversationMember) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(c).Error; err != nil {
			return err
		}
		if len(members) > 0 {
			if err := tx.Create(members).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *ConversationModel) FindByID(id string) (*Conversation, error) {
	var c Conversation
	if err := m.db.Where("id = ?", id).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (m *ConversationModel) UpdateLastActiveAt(id string, t time.Time) error {
	return m.db.Model(&Conversation{}).Where("id = ?", id).Update("last_active_at", t).Error
}

func (m *ConversationModel) CountMembers(conversationID string) (int64, error) {
	var n int64
	err := m.db.Model(&ConversationMember{}).Where("conversation_id = ? AND status = ?", conversationID, "active").Count(&n).Error
	return n, err
}

func (m *ConversationModel) ListByUserID(userID string, offset, limit int) ([]*Conversation, error) {
	var list []*Conversation
	err := m.db.Table("conversations").
		Joins("INNER JOIN conversation_members ON conversation_members.conversation_id = conversations.id AND conversation_members.user_id = ? AND conversation_members.status = ?", userID, "active").
		Order("conversations.last_active_at DESC, conversations.id DESC").
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, err
}

func (m *ConversationModel) ListByUserIDCount(userID string) (int64, error) {
	var n int64
	err := m.db.Model(&ConversationMember{}).Where("user_id = ? AND status = ?", userID, "active").Count(&n).Error
	return n, err
}

func (m *ConversationModel) AddMember(member *ConversationMember) error {
	return m.db.Create(member).Error
}

func (m *ConversationModel) GetMember(conversationID, userID string) (*ConversationMember, error) {
	var member ConversationMember
	if err := m.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (m *ConversationModel) RemoveMember(conversationID, userID string) error {
	return m.db.Model(&ConversationMember{}).
		Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Update("status", "left").Error
}

func (m *ConversationModel) ListMembers(conversationID string) ([]*ConversationMember, error) {
	var list []*ConversationMember
	err := m.db.Where("conversation_id = ?", conversationID).Order("joined_at ASC").Find(&list).Error
	return list, err
}

// FindSingleByTwoUsers 查找 type=single 且仅含 userID1、userID2 两名 active 成员的会话，若不存在返回 nil
func (m *ConversationModel) FindSingleByTwoUsers(userID1, userID2 string) (*Conversation, error) {
	var c Conversation
	err := m.db.Table("conversations").
		Joins("INNER JOIN conversation_members m1 ON m1.conversation_id = conversations.id AND m1.user_id = ? AND m1.status = ?", userID1, "active").
		Joins("INNER JOIN conversation_members m2 ON m2.conversation_id = conversations.id AND m2.user_id = ? AND m2.status = ?", userID2, "active").
		Where("conversations.type = ?", "single").
		First(&c).Error
	if err != nil {
		return nil, err
	}
	n, err := m.CountMembers(c.ID)
	if err != nil {
		return nil, err
	}
	if n != 2 {
		return nil, gorm.ErrRecordNotFound
	}
	return &c, nil
}
