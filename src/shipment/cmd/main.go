package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"shipment/cmd/handlers"
	"shipment/internal/config"
	"shipment/internal/platform"
	"shipment/internal/ship"
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

	err = platform.Migrate(connect, &ship.Shipment{}, &ship.Item{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	repository := ship.NewRepository(connect)

	handler := handlers.NewMessageHandler(repository)

	// Routes
	app.Post("/send-items/:orderId", handler.SendItems)
	app.Get("/shipment/:orderId", handler.FindShipment)
	app.Delete("/shipment/:id", handler.DeleteShipment)

	if err = app.Listen(":9093"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
