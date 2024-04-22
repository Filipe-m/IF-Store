package handlers

import (
	"account/internal/user"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	repository user.Repository
}

func (a *User) CreateUser(c *fiber.Ctx) error {

	var request user.User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := a.repository.Create(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (a *User) UpdateUser(c *fiber.Ctx) error {
	var request user.User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	request.ID = c.Params("id")

	err := a.repository.Update(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (a *User) FindUser(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := a.repository.FindById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (a *User) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	err := a.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func NewUserHandler(repository user.Repository) *User {
	return &User{repository: repository}
}
