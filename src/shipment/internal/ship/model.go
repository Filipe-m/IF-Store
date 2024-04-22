package ship

import (
	"gorm.io/gorm"
	"time"
)

type Shipment struct {
	ID        string         `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	OrderId   string         `json:"order_id" gorm:"column:order_id;type:uuid"`
	Status    string         `json:"status" gorm:"column:status;type:varchar(255)"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
	// one to many
	Items []Item `json:"items" gorm:"foreignKey:ShipmentID;references:ID"`
}

type Item struct {
	ID         string         `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	ShipmentID string         `json:"shipment_id" gorm:"column:shipment_id;type:uuid"`
	ProductID  string         `json:"product_id" gorm:"column:product_id;type:uuid"`
	Quantity   int            `json:"quantity" gorm:"column:quantity;type:integer"`
	CreatedAt  time.Time      `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"column:deleted_at;type:timestamp"`
	// many to one
	Shipment Shipment `json:"-" gorm:"foreignKey:ShipmentID"`
}
