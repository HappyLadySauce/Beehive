package v1

import (
	"time"
)

type Group struct {
	ID          string    `json:"id" 					gorm:"column:id;primaryKey"`	// 群组ID
	Name        string    `json:"name" 					gorm:"column:name;not null" 		validate:"required,min=3,max=32"`	// 群组名称
	Avatar      *string   `json:"avatar,omitempty" 		gorm:"column:avatar;index"`			// 群组头像
	Description *string   `json:"description,omitempty" gorm:"column:description;index"`	// 群组描述
	OwnerID     string    `json:"owner_id" 				gorm:"column:owner_id;not null" 	validate:"required"`	// 群主ID
	Status      string    `json:"status" 				gorm:"column:status;not null" 		validate:"required,oneof=active inactive"`	// 状态	活跃、禁用
	CreatedAt   time.Time `json:"created_at" 			gorm:"column:created_at"`	// 创建时间
	UpdatedAt   time.Time `json:"updated_at" 			gorm:"column:updated_at"`	// 更新时间
}

// TableName 指定表名
func (Group) TableName() string {
	return "groups"
}