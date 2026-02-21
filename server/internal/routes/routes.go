package routes

import (
	"tenderness/internal/handlers"
	"tenderness/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, healthHandler *handlers.HealthHandler, productHandler *handlers.ProductHandler, authHandler *handlers.AuthHandler, oauth2Handler *handlers.OAuth2Handler, jwt *middleware.JWTMiddleware) {
	app.Get("/health", healthHandler.Check)

	api := app.Group("/api")

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/logout", authHandler.Logout)

	// OAuth2 routes
	oauth2 := api.Group("/oauth2")
	oauth2.Get("/:provider/auth", oauth2Handler.GetAuthURL)
	oauth2.Get("/:provider/callback", oauth2Handler.Callback)

	// Protected routes
	protected := api.Group("/user")
	protected.Use(jwt.JWTAuth())
	protected.Get("/profile", authHandler.GetProfile)
	protected.Put("/profile", authHandler.UpdateProfile)
	protected.Put("/password", authHandler.ChangePassword)
	protected.Delete("/account", authHandler.DeleteAccount)

	// OAuth2 linking (protected)
	protected.Post("/link/:provider", oauth2Handler.LinkAccount)
	protected.Delete("/unlink/:provider", oauth2Handler.UnlinkAccount)

	// Product routes
	api.Get("/products", productHandler.GetProducts)
	api.Get("/products/featured", productHandler.GetFeaturedProducts)
	api.Get("/products/search", productHandler.SearchProducts)
	api.Get("/products/category/:category", productHandler.GetProductsByCategory)
	api.Get("/products/:id", productHandler.GetProductByID)

	api.Get("/categories", productHandler.GetCategories)
}
