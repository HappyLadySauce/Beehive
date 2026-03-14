package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GroupJoinRequest 对应 group_join_requests 表
type GroupJoinRequest struct {
	ID             string     `gorm:"column:id;type:uuid;primaryKey"`
	ConversationID string    `gorm:"column:conversation_id;type:varchar(36);not null;uniqueIndex:uq_group_join_conv_user"`
	UserID         string    `gorm:"column:user_id;type:varchar(10);not null;uniqueIndex:uq_group_join_conv_user"`
	Message        string    `gorm:"column:message;type:text;not null;default:''"`
	Status         string    `gorm:"column:status;type:text;not null;default:pending"`
	ProcessedAt    *time.Time `gorm:"column:processed_at;type:timestamptz"`
	ProcessedBy    string    `gorm:"column:processed_by;type:varchar(10)"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (GroupJoinRequest) TableName() string {
	return "group_join_requests"
}

type GroupJoinRequestModel struct {
	db *gorm.DB
}

func NewGroupJoinRequestModel(db *gorm.DB) *GroupJoinRequestModel {
	return &GroupJoinRequestModel{db: db}
}

func (m *GroupJoinRequestModel) Apply(conversationID, userID, message string) (*GroupJoinRequest, error) {
	var existing GroupJoinRequest
	err := m.db.Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&existing).Error
	if err == nil {
		if existing.Status == "pending" {
			return &existing, nil
		}
		if existing.Status == "approved" {
			return nil, gorm.ErrRecordNotFound
		}
		existing.Status = "pending"
		existing.Message = message
		existing.ProcessedAt = nil
		existing.ProcessedBy = ""
		if err := m.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	req := &GroupJoinRequest{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		UserID:         userID,
		Message:        message,
		Status:         "pending",
		CreatedAt:      time.Now(),
	}
	if err := m.db.Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (m *GroupJoinRequestModel) ListPending(conversationID string) ([]*GroupJoinRequest, error) {
	var list []*GroupJoinRequest
	err := m.db.Where("conversation_id = ? AND status = ?", conversationID, "pending").Order("created_at DESC").Find(&list).Error
	return list, err
}

func (m *GroupJoinRequestModel) FindByID(requestID string) (*GroupJoinRequest, error) {
	var req GroupJoinRequest
	if err := m.db.Where("id = ?", requestID).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (m *GroupJoinRequestModel) Approve(requestID, conversationID, processedBy string) error {
	now := time.Now()
	res := m.db.Model(&GroupJoinRequest{}).Where("id = ? AND conversation_id = ? AND status = ?", requestID, conversationID, "pending").Updates(map[string]interface{}{"status": "approved", "processed_at": now, "processed_by": processedBy})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m *GroupJoinRequestModel) Decline(requestID, conversationID, processedBy string) error {
	now := time.Now()
	res := m.db.Model(&GroupJoinRequest{}).Where("id = ? AND conversation_id = ? AND status = ?", requestID, conversationID, "pending").Updates(map[string]interface{}{"status": "declined", "processed_at": now, "processed_by": processedBy})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
