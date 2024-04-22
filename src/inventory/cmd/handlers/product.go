package handlers

import (
	"github.com/gofiber/fiber/v2"
	"inventory/internal/product"
	"inventory/internal/stock"
)

type Product struct {
	repository      product.Repository
	stockRepository stock.Repository
}

func (p *Product) RegisterProduct(c *fiber.Ctx) error {

	var request product.Product
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := p.repository.Create(c.Context(), &request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	err = p.stockRepository.Create(c.Context(), &stock.Stock{
		ProductID: request.ID,
		Quantity:  0,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (p *Product) UpdateProduct(c *fiber.Ctx) error {
	var request product.Product
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

func (p *Product) FindProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	response, err := p.repository.FindById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

func (p *Product) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	err := p.repository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	err = p.stockRepository.Delete(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "success"})
}

func NewProductHandler(repository product.Repository, stockRepository stock.Repository) *Product {
	return &Product{repository: repository, stockRepository: stockRepository}
}
