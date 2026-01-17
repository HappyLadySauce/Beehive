package v1

import (
	"time"
)

type GroupMember struct {
	ID        string    `json:"id" 					gorm:"column:id;primaryKey"`	// 群组成员ID
	GroupID   string    `json:"group_id" 			gorm:"column:group_id;not null" validate:"required"`	// 群组ID
	UserID    string    `json:"user_id" 			gorm:"column:user_id;not null" 	validate:"required"`	// 用户ID
	Role      string    `json:"role" 				gorm:"column:role;not null" 	validate:"required,oneof=owner admin member"`	// 角色	群主、管理员、成员
	Status    string    `json:"status" 				gorm:"column:status;not null" 	validate:"required,oneof=active inactive"`	// 状态	活跃、禁言
	CreatedAt time.Time `json:"created_at" 			gorm:"column:created_at"`	// 创建时间
	UpdatedAt time.Time `json:"updated_at" 			gorm:"column:updated_at"`	// 更新时间
}