package stock

import (
	"gorm.io/gorm"
	"inventory/internal/product"
	"time"
)

type Stock struct {
	ID        string         `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	ProductID string         `json:"product_id" gorm:"column:product_id;type:uuid"`
	Quantity  int            `json:"quantity" gorm:"column:quantity;type:integer"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
	// one to one
	Product product.Product `json:"-" gorm:"foreignKey:ProductID"`
}
