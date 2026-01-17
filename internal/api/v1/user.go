package v1

import (
	"time"
)

type User struct {
	ID           string    `json:"id" 						gorm:"column:id;primaryKey"`
	Nickname     string    `json:"nickname" 				gorm:"column:nickname; not null"				validate:"required,min=3,max=32"`	// 昵称
	Avatar       string    `json:"avatar" 					gorm:"column:avatar; not null"					validate:"required,url,max=255"`	// 头像
	Email        string    `json:"email" 					gorm:"column:email; uniqueIndex; not null"		validate:"required,email,max=255"`	// 邮箱
	Description  string    `json:"description" 				gorm:"column:description; not null"				validate:"required,min=3,max=255"`	// 描述
	Salt         string    `json:"salt,omitempty"			gorm:"column:salt; not null"`			// 盐值
	PasswordHash string    `json:"password_hash,omitempty"	gorm:"column:password_hash; not null"` 	// 密码哈希
	Level        int       `json:"level"					gorm:"column:level; not null"					validate:"required,min=0,max=15"`					// 等级
	Status       string    `json:"status"					gorm:"column:status; not null"					validate:"required,oneof=active inactive"` 	// 状态:活跃、禁用
	CreatedAt    time.Time `json:"created_at"`	// 创建时间
	UpdatedAt    time.Time `json:"updated_at"`	// 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}