package v1

import (
	"time"
)

type Message struct {
	ID          string    `json:"id" 					gorm:"column:id;primaryKey"`	// 消息ID
	Type        string    `json:"type" 					gorm:"column:type;not null" 		validate:"required,oneof=single group"`	// 消息类型	单聊或群聊
	FromUserID  string    `json:"from_user_id" 			gorm:"column:from_user_id;not null" validate:"required"`			// 发送者ID
	ToUserID    *string   `json:"to_user_id,omitempty" 	gorm:"column:to_user_id;index"`										// 接收者ID	单聊时为接收者ID，群聊时为空
	GroupID     *string   `json:"group_id,omitempty" 	gorm:"column:group_id;index"`										// 群组ID	群聊时为群组ID，单聊时为空
	Content     string    `json:"content" 				gorm:"column:content;not null" 		validate:"required"`			// 内容
	MessageType string    `json:"message_type" 			gorm:"column:message_type;not null" validate:"required,oneof=text image file"`	// 消息类型
	Status      string    `json:"status" 				gorm:"column:status;not null" 		validate:"required,oneof=sent delivered read unread"`	// 状态	已发送、已送达、已读、未读
	CreatedAt   time.Time `json:"created_at" 			gorm:"column:created_at"`	// 创建时间
	UpdatedAt   time.Time `json:"updated_at" 			gorm:"column:updated_at"`	// 更新时间
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}