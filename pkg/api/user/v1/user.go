package v1

import (
	"time"
)

// User is the model for the users table.
type User struct {
	ID           string    `json:"id" 						gorm:"column:id;primaryKey"						validate:"required,uuid,eq=10"`	// 用户ID
	Nickname     string    `json:"nickname" 				gorm:"column:nickname; not null"				validate:"required,min=3,max=32"`	// 昵称
	Avatar       string    `json:"avatar" 					gorm:"column:avatar; not null"					validate:"required,url,max=255"`	// 头像
	Email        string    `json:"email" 					gorm:"column:email; uniqueIndex; not null"		validate:"required,email,max=255"`	// 邮箱
	Description  string    `json:"description" 				gorm:"column:description; not null"				validate:"required,min=3,max=255"`	// 描述
	Salt         string    `json:"salt,omitempty"			gorm:"column:salt; not null"`			// 盐值
	PasswordHash string    `json:"password_hash,omitempty"	gorm:"column:password_hash; not null"` 	// 密码哈希
	Level        int       `json:"level"					gorm:"column:level; not null"					validate:"required,min=0,max=15"`	// 等级 0-15 0表示普通用户, 1-15表示账号等级
	Status       string    `json:"status"					gorm:"column:status; not null"					validate:"required,oneof=online offline busy idle invisible"` 	// 状态:在线、离线、忙碌、空闲、隐身
	FreezeTime   time.Time `json:"freeze_time,omitempty"	gorm:"column:freeze_time"						validate:"omitempty,datetime"`	// 冻结时间 冻结后用户无法登录, 为空时表示未冻结
	DeletedAt    time.Time `json:"deleted_at,omitempty"		gorm:"column:deleted_at"						validate:"omitempty,datetime"`	// 删除时间 删除后用户无法登录, 为空时表示未删除
	CreatedAt    time.Time `json:"created_at"				gorm:"column:created_at"						validate:"required,datetime"`	// 创建时间
	UpdatedAt    time.Time `json:"updated_at"				gorm:"column:updated_at"						validate:"required,datetime"`	// 更新时间
}

// TableName returns the table name for the User model.
func (User) TableName() string {
	return "users"
}