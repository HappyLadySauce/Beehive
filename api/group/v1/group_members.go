package v1

import (
	"time"
)

// GroupMember is the model for the group members table.
type GroupMember struct {
	ID        string    `json:"id" 					gorm:"column:id;primaryKey"		validate:"required,uuid,eq=10"`	// 群组成员ID
	GroupID   string    `json:"group_id" 			gorm:"column:group_id;not null" validate:"required,uuid,eq=9"`	// 群组ID
	UserID    string    `json:"user_id" 			gorm:"column:user_id;not null" 	validate:"required,uuid,eq=10"`	// 用户ID
	Role      string    `json:"role" 				gorm:"column:role;not null" 	validate:"required,oneof=owner admin member"`	// 角色	群主、管理员、成员
	InactiveAt time.Time `json:"inactive_at,omitempty"	gorm:"column:inactive_at"		validate:"omitempty,datetime"`	// 禁言时间 禁言后无法发言, 为空时表示未禁言
	DeletedAt time.Time `json:"deleted_at,omitempty"	gorm:"column:deleted_at"		validate:"omitempty,datetime"`	// 删除时间 删除后群组成员无法使用, 为空时表示未删除
	CreatedAt time.Time `json:"created_at" 			gorm:"column:created_at"		validate:"required,datetime"`	// 创建时间
	UpdatedAt time.Time `json:"updated_at" 			gorm:"column:updated_at"		validate:"required,datetime"`	// 更新时间
}

// TableName returns the table name for the GroupMember model.
func (GroupMember) TableName() string {
	return "group_members"
}