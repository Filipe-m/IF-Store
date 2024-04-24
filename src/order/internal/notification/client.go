package notification

import (
	"context"
	"order/internal/httputils/request"
)

type Client interface {
	SendMessage(ctx context.Context, message Message) error
}

type client struct {
	req request.Client
}

func (c *client) SendMessage(ctx context.Context, message Message) error {
	err := c.req.Post(ctx, "/send-message", request.WithRequest(message))
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
