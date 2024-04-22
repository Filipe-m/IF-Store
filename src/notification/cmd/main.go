package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	"notification/cmd/handlers"
	"notification/internal/config"
	"notification/internal/message"
	"notification/internal/platform"
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

	err = platform.Migrate(connect, &message.Message{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	repository := message.NewRepository(connect)

	handler := handlers.NewMessageHandler(repository)

	// Routes
	app.Post("/send-message", handler.SendMessage)

	if err = app.Listen(":9092"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
