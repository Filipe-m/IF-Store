package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string         `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	Username  string         `json:"username" gorm:"column:username;type:varchar(255)"`
	Email     string         `json:"email" gorm:"column:email;type:varchar(255)"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
}
