package product

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          string          `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	Name        string          `json:"name" gorm:"column:name;type:varchar(255)"`
	Description string          `json:"description" gorm:"column:description;type:text"`
	Price       decimal.Decimal `json:"price" gorm:"column:price;type:decimal"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;type:timestamp"`
}
