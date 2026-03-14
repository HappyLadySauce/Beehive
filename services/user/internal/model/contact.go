package model

import (
	"time"

	"gorm.io/gorm"
)

// Contact 对应 contacts 表，owner_id 与 contact_user_id 均为 10 位用户 ID
type Contact struct {
	OwnerID       string    `gorm:"column:owner_id;type:char(10);primaryKey"`
	ContactUserID string    `gorm:"column:contact_user_id;type:char(10);primaryKey"`
	Status        string    `gorm:"column:status;type:text;not null;default:accepted"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamptz;not null"`
}

func (Contact) TableName() string {
	return "contacts"
}

type ContactModel struct {
	db *gorm.DB
}

func NewContactModel(db *gorm.DB) *ContactModel {
	return &ContactModel{db: db}
}

func (m *ContactModel) Add(ownerID, contactUserID string) error {
	return m.db.Create(&Contact{
		OwnerID:       ownerID,
		ContactUserID: contactUserID,
		Status:        "accepted",
		CreatedAt:     time.Now(),
	}).Error
}

func (m *ContactModel) Remove(ownerID, contactUserID string) error {
	return m.db.Where("owner_id = ? AND contact_user_id = ?", ownerID, contactUserID).Delete(&Contact{}).Error
}

func (m *ContactModel) ListContactUserIDs(ownerID string) ([]string, error) {
	var ids []string
	err := m.db.Model(&Contact{}).Where("owner_id = ?", ownerID).Pluck("contact_user_id", &ids).Error
	return ids, err
}

func (m *ContactModel) Exists(ownerID, contactUserID string) (bool, error) {
	var n int64
	err := m.db.Model(&Contact{}).Where("owner_id = ? AND contact_user_id = ?", ownerID, contactUserID).Count(&n).Error
	return n > 0, err
}
