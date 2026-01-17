package v1

import (
	"time"
)

// Message is the model for the messages table.
type Message struct {
	ID          string    `json:"id" 					gorm:"column:id;primaryKey"			validate:"required,uuid,eq=12"`	// 消息ID
	Type        string    `json:"type" 					gorm:"column:type;not null" 		validate:"required,oneof=single group"`	// 消息类型	单聊或群聊
	FromUserID  string    `json:"from_user_id" 			gorm:"column:from_user_id;not null" validate:"required,uuid,eq=10"`		// 发送者ID
	ToUserID    *string   `json:"to_user_id,omitempty" 	gorm:"column:to_user_id;index"		validate:"omitempty,uuid,eq=10"`	// 接收者ID	单聊时为接收者ID，群聊时为空
	GroupID     *string   `json:"group_id,omitempty" 	gorm:"column:group_id;index"		validate:"omitempty,uuid,eq=9"`	// 群组ID	群聊时为群组ID，单聊时为空
	Content     string    `json:"content" 				gorm:"column:content;not null" 		validate:"required,min=1,max=4096"`	// 内容 1-4096字符
	MessageType string    `json:"message_type" 			gorm:"column:message_type;not null" validate:"required,oneof=text image file"`	// 消息类型 文本、图片、文件
	Status      string    `json:"status" 				gorm:"column:status;not null" 		validate:"required,oneof=sent delivered read unread"`	// 状态	已发送、已送达、已读、未读
	DeletedAt   time.Time `json:"deleted_at,omitempty"	gorm:"column:deleted_at"			validate:"omitempty,datetime"`	// 删除时间 删除后消息无法查看, 为空时表示未删除
	CreatedAt   time.Time `json:"created_at" 			gorm:"column:created_at"			validate:"required,datetime"`	// 创建时间
	UpdatedAt   time.Time `json:"updated_at" 			gorm:"column:updated_at"			validate:"required,datetime"`	// 更新时间
}

// TableName returns the table name for the Message model.
func (Message) TableName() string {
	return "messages"
}