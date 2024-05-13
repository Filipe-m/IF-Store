package mail

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"notification/internal/chat"
	"notification/internal/config"
	"time"
)

type client struct {
	cfg config.Mail
}

func (c *client) SendMessage(ctx context.Context, to string, message chat.Message) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	if c.cfg.Username == "" || c.cfg.Password == "" {
		log.Println("mail: username or password is empty")
		return nil
	}

	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s", to, message.Body))

	return smtp.SendMail(fmt.Sprintf("%s:%s", c.cfg.Host, c.cfg.Port), auth, c.cfg.Username, []string{to}, msg)
}

func NewClient(cfg config.Mail) chat.Client {
	return &client{cfg: cfg}
}
