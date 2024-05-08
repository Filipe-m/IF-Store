package order

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	PENDING    Status = "pending"
	PROCESSING Status = "processing"
	FINISHED   Status = "finished"
	CANCELED   Status = "cancel"
)

type Order struct {
	ID          string          `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	UserID      string          `json:"user_id" gorm:"column:user_id;type:uuid"`
	Status      Status          `json:"status" gorm:"column:status;type:varchar(255);default:pending"`
	TotalAmount decimal.Decimal `json:"total_amount" gorm:"column:total_amount;type:decimal"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;type:timestamp"`
	// one to many
	Items []Item `json:"items" gorm:"foreignKey:OrderID;references:ID"`
}

type Item struct {
	ID        string          `json:"id" gorm:"column:id;primarykey;type:uuid;default:gen_random_uuid()"`
	OrderID   string          `json:"order_id" gorm:"column:order_id;type:uuid"`
	ProductID string          `json:"product_id" gorm:"column:product_id;type:uuid"`
	Quantity  int             `json:"quantity" gorm:"column:quantity;type:integer"`
	UnitPrice decimal.Decimal `json:"unit_price" gorm:"column:unit_price;type:decimal"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	DeletedAt gorm.DeletedAt  `json:"-" gorm:"column:deleted_at;type:timestamp"`
	// many to one
	Order Order `json:"-" gorm:"foreignKey:OrderID"`
}

type FinishOrder struct {
	PaymentData string `json:"payment_data"`
	OrderID     string `json:"order_id"`
}

func (i *Item) CalculateTotalPrice() decimal.Decimal {
	return i.UnitPrice.Mul(decimal.NewFromInt(int64(i.Quantity)))
}

func (o *Order) ProductIDs() []string {
	products := make([]string, len(o.Items))
	for i, item := range o.Items {
		products[i] = item.ProductID
	}
	return products
}
