package handlers

import (
	"github.com/gofiber/fiber/v2"
	"inventory/internal/stock"
)

type Stock struct {
	repository stock.Repository
}

func (p *Stock) UpdateStock(c *fiber.Ctx) error {
	var request stock.Stock
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	request.ID = c.Params("id")

	err := p.repository.Update(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (p *Stock) FindStock(c *fiber.Ctx) error {
	id := c.Params("productId")

	response, err := p.repository.FindByProductId(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func NewStockHandler(repository stock.Repository) *Stock {
	return &Stock{repository: repository}
}
