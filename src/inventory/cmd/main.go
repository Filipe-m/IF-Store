package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
		log.Printf("[LOAD ENVIRONMENT VARIABLES FAIL]: %s\n", err.Error())
	}
}

func main() {
	app := fiber.New()
	app.Use(logger.New())

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
	app.Use(cors.New())
	app.Post("/product/register", productHandler.RegisterProduct)
	app.Put("/product/:id", productHandler.UpdateProduct)
	app.Get("/product/:id", productHandler.FindProduct)
	app.Delete("/product/:id", productHandler.DeleteProduct)

	app.Put("/stock/:productId/add", stockHandler.AddStock)
	app.Put("/stock/:productId/remove", stockHandler.RemoveStock)
	app.Get("/stock/:productId", stockHandler.FindStock)

	app.Get("/product", productHandler.ListProducts)

	if err = app.Listen(":9094"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
