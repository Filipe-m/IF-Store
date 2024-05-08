package inventory

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	"order/internal/httputils/request"
)

type Client interface {
	RemoveItemsFromStock(ctx context.Context, products []Product) error
	CheckStock(ctx context.Context, products []Product) error
	FindItemPrice(ctx context.Context, productId string) (decimal.Decimal, error)
	GetProducts(ctx context.Context, productIds []string) ([]Product, error)
}

type client struct {
	req request.Client
}

func (c *client) RemoveItemsFromStock(ctx context.Context, products []Product) error {
	for _, product := range products {
		err := c.req.Put(ctx, fmt.Sprintf("/stock/%s/remove", product.ID), request.WithRequest(Stock{Quantity: product.Quantity}))
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *client) GetProducts(ctx context.Context, productIds []string) ([]Product, error) {
	products := make([]Product, len(productIds))
	for i, id := range productIds {
		err := c.req.Get(ctx, fmt.Sprintf("/product/%s", id), request.WithResponse(&products[i]))
		if err != nil {
			return nil, err
		}
		products[i].ID = id
	}
	return products, nil
}

func (c *client) CheckStock(ctx context.Context, products []Product) error {
	for _, product := range products {
		var stock Stock
		err := c.req.Get(ctx, fmt.Sprintf("/stock/%s", product.ID), request.WithResponse(&stock))
		if err != nil {
			return err
		}

		if stock.Quantity < product.Quantity {
			return fmt.Errorf("insufficient stock for product %s", product.ID)
		}
	}
	return nil
}

func (c *client) FindItemPrice(ctx context.Context, productId string) (decimal.Decimal, error) {
	var prd Product
	err := c.req.Get(ctx, fmt.Sprintf("/product/%s", productId), request.WithResponse(&prd))
	if err != nil {
		return decimal.Decimal{}, err
	}

	return prd.Price, nil
}

func NewClient(req request.Client) Client {
	return &client{
		req: req,
	}
}
