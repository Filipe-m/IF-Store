package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"order/internal/inventory"
	"order/internal/order"
)

type Order struct {
	repository order.Repository
	inventory  inventory.Client
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

func (a *Order) RemoveItemToOrder(c *fiber.Ctx) error {

	userID := c.Get("USER-ID")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	var request order.Item
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	resp, err := a.service.RemoveItemToOrder(c.Context(), userID, &request)
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

func (a *Order) FindActualOrder(c *fiber.Ctx) error {
	userID := c.Get("USER-ID")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthorized"})
	}

	response, err := a.repository.FindActualByUserId(c.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "not found"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	products, err := a.inventory.GetProducts(c.Context(), response.ProductIDs())
	if err != nil {
		return err
	}

	for i, item := range response.Items {
		for _, product := range products {
			if item.ProductID == product.ID {
				response.Items[i].Name = product.Name
			}
		}
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

func NewOrderHandler(repository order.Repository, invetory inventory.Client, service order.Service) *Order {
	return &Order{repository: repository, inventory: invetory, service: service}
}
