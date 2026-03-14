package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContactRequest 对应 contact_requests 表
type ContactRequest struct {
	ID         string    `gorm:"column:id;type:uuid;primaryKey"`
	FromUserID string    `gorm:"column:from_user_id;type:varchar(10);not null;uniqueIndex:uq_contact_request_pair"`
	ToUserID   string    `gorm:"column:to_user_id;type:varchar(10);not null;uniqueIndex:uq_contact_request_pair"`
	Status     string    `gorm:"column:status;type:text;not null;default:pending"`
	Message    string    `gorm:"column:message;type:text;not null;default:''"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamptz;not null"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamptz;not null"`
}

func (ContactRequest) TableName() string {
	return "contact_requests"
}

type ContactRequestModel struct {
	db *gorm.DB
}

func NewContactRequestModel(db *gorm.DB) *ContactRequestModel {
	return &ContactRequestModel{db: db}
}

// CreateOrReapply 创建申请；若已存在且为 declined 则更新为 pending 并返回
func (m *ContactRequestModel) CreateOrReapply(fromUserID, toUserID, message string) (*ContactRequest, error) {
	now := time.Now()
	var existing ContactRequest
	err := m.db.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).First(&existing).Error
	if err == nil {
		if existing.Status == "pending" {
			return &existing, nil
		}
		if existing.Status == "accepted" {
			return nil, gorm.ErrRecordNotFound // 已通过，可返回已存在或业务错误
		}
		// declined：更新为 pending
		existing.Status = "pending"
		existing.Message = message
		existing.UpdatedAt = now
		if err := m.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	req := &ContactRequest{
		ID:         uuid.New().String(),
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     "pending",
		Message:    message,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := m.db.Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

func (m *ContactRequestModel) ListPendingByToUser(toUserID string) ([]*ContactRequest, error) {
	var list []*ContactRequest
	err := m.db.Where("to_user_id = ? AND status = ?", toUserID, "pending").Order("created_at DESC").Find(&list).Error
	return list, err
}

func (m *ContactRequestModel) FindByID(requestID string) (*ContactRequest, error) {
	var req ContactRequest
	if err := m.db.Where("id = ?", requestID).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (m *ContactRequestModel) Accept(requestID, toUserID string) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		var req ContactRequest
		if err := tx.Where("id = ? AND to_user_id = ? AND status = ?", requestID, toUserID, "pending").First(&req).Error; err != nil {
			return err
		}
		now := time.Now()
		if err := tx.Model(&ContactRequest{}).Where("id = ?", requestID).Updates(map[string]interface{}{"status": "accepted", "updated_at": now}).Error; err != nil {
			return err
		}
		// 双向添加好友
		if err := tx.Create(&Contact{OwnerID: toUserID, ContactUserID: req.FromUserID, Status: "accepted", CreatedAt: now}).Error; err != nil {
			return err
		}
		if err := tx.Create(&Contact{OwnerID: req.FromUserID, ContactUserID: toUserID, Status: "accepted", CreatedAt: now}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (m *ContactRequestModel) Decline(requestID, toUserID string) error {
	res := m.db.Model(&ContactRequest{}).Where("id = ? AND to_user_id = ? AND status = ?", requestID, toUserID, "pending").Updates(map[string]interface{}{"status": "declined", "updated_at": time.Now()})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
