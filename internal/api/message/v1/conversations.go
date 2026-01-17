package v1

import (
	"time"
)

// Conversation is the model for the conversations table.
type Conversation struct {
	ID        string    	`json:"id" 							gorm:"column:id;primaryKey"			validate:"required,uuid,eq=11"`	// 会话ID
	Type      string    	`json:"type" 						gorm:"column:type;not null" 		validate:"required,oneof=single group"`	// 会话类型	单聊或群聊
	UserID1   *string   	`json:"user_id1,omitempty" 			gorm:"column:user_id1;index"		validate:"omitempty,uuid,eq=10"`	// 用户1ID	单聊时为用户1ID，群聊时为空
	UserID2   *string   	`json:"user_id2,omitempty" 			gorm:"column:user_id2;index"		validate:"omitempty,uuid,eq=10"`	// 用户2ID	单聊时为用户2ID，群聊时为空
	GroupID   *string   	`json:"group_id,omitempty" 			gorm:"column:group_id;index"		validate:"omitempty,uuid,eq=9"`	// 群组ID	群聊时为群组ID，单聊时为空
	LastMessageID *string   `json:"last_message_id,omitempty" 	gorm:"column:last_message_id;index"	validate:"omitempty,uuid,eq=12"`	// 最后消息ID	最后消息ID
	LastMessageAt *time.Time `json:"last_message_at,omitempty" 	gorm:"column:last_message_at;index"	validate:"omitempty,datetime"`	// 最后消息时间	最后消息时间
	UnreadCount int    		`json:"unread_count" 				gorm:"column:unread_count;not null;default:0"	validate:"required,min=0,max=999"`	// 未读消息数 0-999
	DeletedAt   time.Time 	`json:"deleted_at,omitempty"		gorm:"column:deleted_at"			validate:"omitempty,datetime"`	// 删除时间 删除后会话无法使用, 为空时表示未删除
	CreatedAt time.Time 	`json:"created_at" 					gorm:"column:created_at"			validate:"required,datetime"`	// 创建时间
	UpdatedAt time.Time 	`json:"updated_at" 					gorm:"column:updated_at"			validate:"required,datetime"`	// 更新时间
}

// TableName returns the table name for the Conversation model.
func (Conversation) TableName() string {
	return "conversations"
}