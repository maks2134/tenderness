package routes

import (
	"tenderness/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler) {
	api := app.Group("/api")
	app.Get("/health", healthHandler.Check)
	api.Get("/users", func(c *fiber.Ctx) error { return c.SendString("Users") })
}
