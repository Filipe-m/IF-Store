package handlers

import (
	"github.com/gofiber/fiber/v2"
	"order/internal/order"
)

type Order struct {
	repository order.Repository
	service    order.Service
}

func (a *Order) AddItemToOrder(c *fiber.Ctx) error {

	userID := c.Get("USER-ID")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	var request order.Item
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	resp, err := a.service.AddItemToOrder(c.Context(), userID, &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (a *Order) FinishOrder(c *fiber.Ctx) error {

	var request order.FinishOrder
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := a.service.FinishOrder(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func (a *Order) FindOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := a.repository.FindByOrderId(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (a *Order) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")

	err := a.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func NewOrderHandler(repository order.Repository, service order.Service) *Order {
	return &Order{repository: repository, service: service}
}
