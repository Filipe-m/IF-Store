package message

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID        string         `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	UserID    string         `json:"user_id" gorm:"column:user_id;type:uuid"`
	OrderId   string         `json:"order_id" gorm:"column:order_id;type:uuid"`
	Message   string         `json:"message" gorm:"column:message;type:text"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
}
