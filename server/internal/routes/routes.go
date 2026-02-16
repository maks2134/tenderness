package routes

import (
	"tenderness/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler, productHandler *handlers.ProductHandler) {
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")

	api.Get("/products", productHandler.GetProducts)
	api.Get("/products/featured", productHandler.GetFeaturedProducts)
	api.Get("/products/search", productHandler.SearchProducts)
	api.Get("/products/category/:category", productHandler.GetProductsByCategory)
	api.Get("/products/:id", productHandler.GetProductByID)

	api.Get("/categories", productHandler.GetCategories)
}
