package handlers

import (
	"github.com/gofiber/fiber/v2"
	"notification/internal/message"
)

type User struct {
	repository message.Repository
}

func (a *User) SendMessage(c *fiber.Ctx) error {

	var request message.Message
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := a.repository.Create(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func NewMessageHandler(repository message.Repository) *User {
	return &User{repository: repository}
}
