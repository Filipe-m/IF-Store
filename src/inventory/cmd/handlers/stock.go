package handlers

import (
	"github.com/gofiber/fiber/v2"
	"inventory/internal/stock"
)

type Stock struct {
	repository stock.Repository
}

func (p *Stock) RemoveStock(c *fiber.Ctx) error {
	var request stock.Stock
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	request.ProductID = c.Params("productId")

	stockModel, err := p.repository.FindByProductId(c.Context(), request.ProductID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if stockModel.Quantity < request.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Insufficient stock"})
	}

	stockModel.Quantity = stockModel.Quantity - request.Quantity

	err = p.repository.Update(c.Context(), stockModel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(stockModel)
}

func (p *Stock) AddStock(c *fiber.Ctx) error {
	var request stock.Stock
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	request.ProductID = c.Params("productId")

	stockModel, err := p.repository.FindByProductId(c.Context(), request.ProductID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	stockModel.Quantity += request.Quantity

	err = p.repository.Update(c.Context(), stockModel)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(stockModel)
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
