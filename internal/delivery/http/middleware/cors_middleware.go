package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/irfantrue/web-application-backend-go/internal/config"
)

func CorsMiddleware(cfg *config.Config) fiber.Handler {
	origins := cfg.App.CorsAllowedOrigins
	allowCredentials := true
	if origins == "" {
		origins = "*"
		allowCredentials = false
	}

	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ", "),
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Timezone"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
		AllowCredentials: allowCredentials,
	})
}
