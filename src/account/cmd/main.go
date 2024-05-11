package main

import (
	"account/cmd/handlers"
	"account/internal/config"
	"account/internal/platform"
	"account/internal/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
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

	err = platform.Migrate(connect, user.User{})
	if err != nil {
		log.Fatalf("[MIGRATE DATABASE FAIL]: %s", err.Error())
	}

	repository := user.NewRepository(connect)

	handler := handlers.NewUserHandler(repository)

	// Routes
	app.Use(cors.New())
	app.Post("/users", handler.CreateUser)
	app.Put("/users/:id", handler.UpdateUser)
	app.Get("/users/:id", handler.FindUser)
	app.Delete("/users/:id", handler.DeleteUser)

	if err = app.Listen(":9091"); err != nil {
		log.Fatalf("[START SERVER FAIL]: %s", err.Error())
	}
}
