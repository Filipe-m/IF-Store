package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"notification/cmd/handlers"
	"notification/internal/account"
	mail2 "notification/internal/chat/mail"
	"notification/internal/config"
	"notification/internal/httputils/request"
	"notification/internal/message"
	"notification/internal/platform"
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

	err = platform.Migrate(connect, &message.Message{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	repository := message.NewRepository(connect)

	accountClient := account.NewClient(request.New(cfg.Service.AccountURL))

	mail := mail2.NewClient(cfg.Mail)

	service := message.NewService(repository, mail, accountClient)

	handler := handlers.NewMessageHandler(service)

	// Routes
	app.Post("/send-message", handler.SendMessage)

	if err = app.Listen(":9092"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
