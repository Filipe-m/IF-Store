package payment

import (
	"context"
	"order/internal/httputils/request"
)

type Client interface {
	SendPayment(context.Context) error
}

type client struct {
	req request.Client
}

func (c *client) SendPayment(ctx context.Context) error {
	return nil
}

func NewClient(req request.Client) Client {
	return &client{
		req: req,
	}
}
