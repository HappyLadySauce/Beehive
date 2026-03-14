package model

import (
	"time"

	"gorm.io/gorm"
)

// User 对应 users 表，供 AuthService 进行账号密码校验使用。
type User struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey"`
	Username     string    `gorm:"column:username;type:text;uniqueIndex;not null"`
	PasswordHash string    `gorm:"column:password_hash;type:text;not null"`
	Status       string    `gorm:"column:status;type:text;not null;default:'normal'"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamptz;not null;default:now()"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:timestamptz;not null;default:now()"`
}

func (User) TableName() string {
	return "users"
}

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

func (m *UserModel) FindByUsername(username string) (*User, error) {
	var u User
	if err := m.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Create 创建用户，用于注册。调用方需保证 ID、Username、PasswordHash、Status 已填。
func (m *UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}

