package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"order/cmd/handlers"
	"order/internal/config"
	"order/internal/httputils/request"
	"order/internal/inventory"
	"order/internal/notification"
	"order/internal/order"
	"order/internal/payment"
	"order/internal/platform"
	"order/internal/shipment"
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

	err = platform.Migrate(connect, order.Order{}, order.Item{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	repository := order.NewRepository(connect)

	inventory := inventory.NewClient(request.New(cfg.Service.InventoryUrl))
	payment := payment.NewClient(request.New(cfg.Service.PaymentUrl))
	notification := notification.NewClient(request.New(cfg.Service.NotificationUrl))
	shipment := shipment.NewClient(request.New(cfg.Service.ShipmentUrl))

	service := order.NewService(repository, inventory, payment, notification, shipment)

	handler := handlers.NewOrderHandler(repository, service)

	// Routes
	app.Post("/order-item", handler.AddItemToOrder)
	app.Get("/order/:id", handler.FindOrder)
	app.Delete("/order/:id", handler.DeleteOrder)
	app.Post("/order/finish", handler.FinishOrder)

	if err = app.Listen(":9095"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
