// Package http provides HTTP delivery components including route setup and fiber instance initialization.
package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/handler"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/middleware"
)

// Config holds the dependencies required to set up HTTP routes.
type Config struct {
	App            *fiber.App
	Config         *config.Config
	ProductHandler *handler.ProductHandler
}

// NewSetupRoutes configures the global middlewares and registers all API endpoints.
func NewSetupRoutes(cfg *Config) {
	// 1. RECOVER (Must always be number 1 to catch panics)
	cfg.App.Use(middleware.RecoverMiddleware())

	// 2. CORS (Placed second so the browser obtains cross-domain permission quickly)
	cfg.App.Use(middleware.CorsMiddleware(cfg.Config))

	// 3. HELMET (Third in the sequence for injecting Security Headers)
	cfg.App.Use(middleware.HelmetMiddleware())

	// 4. Other Global Middleware (Compression and Timezone)
	cfg.App.Use(middleware.CompressMiddleware(cfg.Config))
	cfg.App.Use(middleware.TimezoneMiddleware("Asia/Jakarta"))

	v1 := cfg.App.Group("/api/v1")

	products := v1.Group("/products")
	{
		products.Post("/", cfg.ProductHandler.Create)
	}
}
