package model

import (
	"time"
)

type User struct {
	ID           string    `json:"id" gorm:"primaryKey"`

	// 昵称
	Nickname     string    `json:"nickname" gorm:"not null" validate:"required,min=3,max=32"`

	// 头像
	Avatar       string    `json:"avatar" gorm:"not null" validate:"required,url,max=255"`

	// 邮箱
	Email        string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email,max=255"`

	// 盐值，如果为空则自动生成
	Salt         string    `json:"salt,omitempty" gorm:"column:salt;not null"`
	// 密码哈希，如果为空则自动生成
	PasswordHash string    `json:"password_hash,omitempty" gorm:"column:password_hash;not null"` 

	Level        int       `json:"level" gorm:"not null" validate:"required,min=1,max=10"`
	
	// 状态
	Status       string    `json:"status" gorm:"not null" validate:"required,oneof=active inactive deleted"`

	CreatedAt    time.Time `json:"created_at"` // 由 GORM 自动设置
	UpdatedAt    time.Time `json:"updated_at"` // 由 GORM 自动设置
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}