package payment

import (
	"context"
	"fmt"
	"order/internal/httputils/request"
)

type Client interface {
	SendPayment(userId string, payment Payment, ctx context.Context) error
}

type client struct {
	req request.Client
}

func (c *client) SendPayment(userId string, payment Payment, ctx context.Context) error {
	err := c.req.Post(ctx, fmt.Sprintf("/payment/%s", userId), request.WithRequest([]Payment{payment}))
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
