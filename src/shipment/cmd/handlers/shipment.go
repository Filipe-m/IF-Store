package handlers

import (
	"github.com/gofiber/fiber/v2"
	"shipment/internal/ship"
)

type Shipment struct {
	repository ship.Repository
}

func (s *Shipment) SendItems(c *fiber.Ctx) error {

	var request []ship.Item
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	id := c.Params("orderId")

	err := s.repository.Create(c.Context(), &ship.Shipment{OrderId: id, Items: request})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (s *Shipment) FindShipment(c *fiber.Ctx) error {
	id := c.Params("orderId")

	response, err := s.repository.FindByOrderId(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (s *Shipment) DeleteShipment(c *fiber.Ctx) error {
	id := c.Params("id")

	err := s.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func NewMessageHandler(repository ship.Repository) *Shipment {
	return &Shipment{repository: repository}
}
