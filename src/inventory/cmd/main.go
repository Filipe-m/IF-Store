package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"inventory/cmd/handlers"
	"inventory/internal/config"
	"inventory/internal/platform"
	"inventory/internal/product"
	"inventory/internal/stock"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("[LOAD ENVIRONMENT VARIABLES FAIL]: %s", err.Error())
	}
}

func main() {
	app := fiber.New()

	cfg := config.Load()

	connect, err := platform.NewPostgresConnect(cfg.Database)
	if err != nil {
		log.Fatalf("[CONNECT DATABASE FAIL]: %s", err.Error())
	}

	err = platform.Migrate(connect, &product.Product{}, &stock.Stock{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	productRepository := product.NewRepository(connect)

	stockRepository := stock.NewRepository(connect)

	productHandler := handlers.NewProductHandler(productRepository, stockRepository)

	stockHandler := handlers.NewStockHandler(stockRepository)

	// Routes
	app.Post("/product/register", productHandler.RegisterProduct)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:id", productHandler.FindProduct)
	app.Delete("/product/:id", productHandler.DeleteProduct)

	app.Put("/stock/:id", stockHandler.UpdateStock)
	app.Get("/stock/:productId", stockHandler.FindStock)

	app.Get("/product", productHandler.ListProducts)

	if err = app.Listen(":9094"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
