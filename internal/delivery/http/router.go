// Package http provides HTTP delivery components including route setup and fiber instance initialization.
package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/handler"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/middleware"
)

// Config holds the dependencies required to set up HTTP routes.
type Config struct {
	App           *fiber.App
	TokenService  port.TokenService
	APIKeyService usecase.APIKeyUsecase
	Config        *config.Config
	APIKeyHandler *handler.APIKeyHandler
	AuthHandler   *handler.AuthHandler
}

// NewSetupRoutes configures the global middlewares and registers all API endpoints.
func NewSetupRoutes(cfg *Config) {
	// 1. RECOVER (Must always be number 1 to catch panics)
	cfg.App.Use(middleware.RecoverMiddleware())

	// 2. CORS (Placed second so the browser obtains cross-domain permission quickly)
	cfg.App.Use(middleware.CorsMiddleware(cfg.Config))

	// 3. HELMET (Third in the sequence for injecting Security Headers)
	cfg.App.Use(middleware.HelmetMiddleware())

	// 4. Other Global Middleware
	cfg.App.Use(middleware.CompressMiddleware(cfg.Config))
	cfg.App.Use(middleware.TimezoneMiddleware("Asia/Jakarta"))

	v1 := cfg.App.Group("/api/v1")

	protectedAPIKey := v1.Group("/")
	protectedAPIKey.Use(middleware.APIKeyMiddleware(cfg.APIKeyService))

	auths := protectedAPIKey.Group("/auths")
	{
		auths.Post("/register", cfg.AuthHandler.Register)
		auths.Post("/login", cfg.AuthHandler.Login)
		auths.Post("/logout", cfg.AuthHandler.Logout)
		auths.Post("/refresh", cfg.AuthHandler.Refresh)
		auths.Post("/forgot-password", cfg.AuthHandler.ForgotPassword)
		auths.Post("/reset-password", cfg.AuthHandler.ResetPassword)
	}

	protectedToken := v1.Group("/")
	protectedToken.Use(middleware.JwtMiddleware(cfg.TokenService))

	apiKeyManagements := protectedToken.Group("/api-key-managements")
	{
		apiKeyManagements.Post("/", cfg.APIKeyHandler.Generate)
		apiKeyManagements.Get("/", cfg.APIKeyHandler.FindAll)
		apiKeyManagements.Patch("/:id/activate", cfg.APIKeyHandler.Activate)
		apiKeyManagements.Patch("/:id/revoke", cfg.APIKeyHandler.Revoke)
	}
}
