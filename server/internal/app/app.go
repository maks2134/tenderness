package app

import (
	"log"
	"os"

	"tenderness/internal/configs"
	"tenderness/internal/domain/storage"
	"tenderness/internal/handlers"
	"tenderness/internal/middleware"
	"tenderness/internal/repository"
	"tenderness/internal/routes"
	"tenderness/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Run() {
	config := configs.LoadConfig()

	db, err := storage.NewDatabase(config.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "/root/migrations"
	}
	if err := db.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Product repositories and services
	productRepo := repository.NewProductRepository(db.DB)
	categoryRepo := repository.NewCategoryRepository(db.DB)
	productService := services.NewProductService(productRepo, categoryRepo)

	// Auth repositories and services
	userRepo := repository.NewUserRepository(db.DB)
	jwtMiddleware := middleware.NewJWTMiddleware("your-secret-key-here") // TODO: Move to config
	authService := services.NewAuthService(userRepo, jwtMiddleware)
	oauth2Service := services.NewOAuth2Service(userRepo, jwtMiddleware)

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	productHandler := handlers.NewProductHandler(productService)
	authHandler := handlers.NewAuthHandler(authService, jwtMiddleware)
	oauth2Handler := handlers.NewOAuth2Handler(oauth2Service)

	app := fiber.New(fiber.Config{
		AppName: "Tenderness App",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	routes.SetupRoutes(app, healthHandler, productHandler, authHandler, oauth2Handler, jwtMiddleware)

	log.Println("Starting server on port " + config.Port)
	log.Fatal(app.Listen(config.Port))
}
