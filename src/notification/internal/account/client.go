package account

import (
	"context"
	"fmt"
	"notification/internal/httputils/request"
)

type Client interface {
	GetUser(ctx context.Context, id string) (*User, error)
}

type client struct {
	req request.Client
}

func (c *client) GetUser(ctx context.Context, id string) (*User, error) {
	var user User
	err := c.req.Get(ctx, fmt.Sprintf("/users/%s", id), request.WithResponse(&user))
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func NewClient(req request.Client) Client {
	return &client{
		req: req,
	}
}
