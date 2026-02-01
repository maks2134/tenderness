package app

import (
	"tenderness/internal/configs"
	"tenderness/internal/handlers"
	"tenderness/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run() {
	config := configs.LoadConfig()
	app := fiber.New(fiber.Config{
		AppName: "Tenderness App",
	})
	app.Use(logger.New())
	healthHandler := handlers.NewHealthHandler()
	routes.SetupRoutes(app, healthHandler)
	log.Fatal(app.Listen(config.Port))
}
