package inventory

import (
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          string          `json:"product_id"`
	Quantity    int             `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
}

type Stock struct {
	Quantity int `json:"quantity"`
}
