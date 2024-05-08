package message

import (
	"context"
	"notification/internal/account"
	"notification/internal/chat"
)

type Service interface {
	SendMessage(ctx context.Context, message *Message) error
}

type service struct {
	repository Repository
	account    account.Client
	mail       chat.Client
}

func (s *service) SendMessage(ctx context.Context, message *Message) error {
	user, err := s.account.GetUser(ctx, message.UserID)
	if err != nil {
		return err
	}

	err = s.mail.SendMessage(ctx, user.Email, chat.Message{
		Subject: message.Message,
		Body:    message.Message,
	})
	if err != nil {
		return err
	}

	err = s.repository.Create(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

func NewService(repository Repository, mail chat.Client, account account.Client) Service {
	return &service{repository: repository, mail: mail, account: account}
}
