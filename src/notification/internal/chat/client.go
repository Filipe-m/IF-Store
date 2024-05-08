package chat

import "context"

type Client interface {
	SendMessage(ctx context.Context, to string, message Message) error
}
