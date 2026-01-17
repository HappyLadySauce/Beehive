package v1

import (
	"time"
)

type Conversation struct {
	ID        string    	`json:"id" 							gorm:"column:id;primaryKey"`	// 会话ID
	Type      string    	`json:"type" 						gorm:"column:type;not null" 		validate:"required,oneof=single group"`	// 会话类型	单聊或群聊
	UserID1   *string   	`json:"user_id1,omitempty" 			gorm:"column:user_id1;index"`	// 用户1ID	单聊时为用户1ID，群聊时为空
	UserID2   *string   	`json:"user_id2,omitempty" 			gorm:"column:user_id2;index"`	// 用户2ID	单聊时为用户2ID，群聊时为空
	GroupID   *string   	`json:"group_id,omitempty" 			gorm:"column:group_id;index"`	// 群组ID	群聊时为群组ID，单聊时为空
	LastMessageID *string   `json:"last_message_id,omitempty" 	gorm:"column:last_message_id;index"`	// 最后消息ID	最后消息ID
	LastMessageAt *time.Time `json:"last_message_at,omitempty" 	gorm:"column:last_message_at;index"`	// 最后消息时间	最后消息时间
	UnreadCount int    		`json:"unread_count" 				gorm:"column:unread_count;not null;default:0"`	// 未读消息数
	Status    string    	`json:"status" 						gorm:"column:status;not null" 		validate:"required,oneof=active inactive"`	// 状态	活跃、禁言
	CreatedAt time.Time 	`json:"created_at" 					gorm:"column:created_at"`	// 创建时间
	UpdatedAt time.Time 	`json:"updated_at" 					gorm:"column:updated_at"`	// 更新时间
}