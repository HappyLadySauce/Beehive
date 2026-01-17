package v1

import (
	"time"
)

// Group is the model for the groups table.
type Group struct {
	ID          string    `json:"id" 					gorm:"column:id;primaryKey"			validate:"required,uuid,eq=9"`	// 群组ID
	Name        string    `json:"name" 					gorm:"column:name;not null" 		validate:"required,min=3,max=32"`	// 群组名称
	Avatar      *string   `json:"avatar,omitempty" 		gorm:"column:avatar;index"			validate:"omitempty,url,max=255"`	// 群组头像
	Description *string   `json:"description,omitempty" gorm:"column:description;index"		validate:"omitempty,min=3,max=255"`	// 群组描述 3-255字符
	OwnerID     string    `json:"owner_id" 				gorm:"column:owner_id;not null" 	validate:"required,uuid,eq=10"`	// 群主ID
	DeletedAt   time.Time `json:"deleted_at,omitempty"	gorm:"column:deleted_at"			validate:"omitempty,datetime"`	// 删除时间 删除后群组无法使用, 为空时表示未删除
	CreatedAt   time.Time `json:"created_at" 			gorm:"column:created_at"			validate:"required,datetime"`	// 创建时间
	UpdatedAt   time.Time `json:"updated_at" 			gorm:"column:updated_at"			validate:"required,datetime"`	// 更新时间
}

// TableName returns the table name for the Group model.
func (Group) TableName() string {
	return "groups"
}