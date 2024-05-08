package handlers

import (
	"github.com/gofiber/fiber/v2"
	"notification/internal/message"
)

type User struct {
	service message.Service
}

func (a *User) SendMessage(c *fiber.Ctx) error {

	var request message.Message
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := a.service.SendMessage(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func NewMessageHandler(service message.Service) *User {
	return &User{service: service}
}
