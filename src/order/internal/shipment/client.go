package shipment

import (
	"context"
	"fmt"
	"order/internal/httputils/request"
)

type Client interface {
	SendItems(ctx context.Context, items []Item, orderId string) error
}

type client struct {
	req request.Client
}

func (c *client) SendItems(ctx context.Context, items []Item, orderId string) error {
	err := c.req.Post(ctx, fmt.Sprintf("/send-items/%s", orderId), request.WithRequest(items))
	if err != nil {
		return err
	}
	return nil
}

func NewClient(req request.Client) Client {
	return &client{
		req: req,
	}
}
