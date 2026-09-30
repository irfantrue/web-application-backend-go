package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/irfantrue/web-application-backend-go/internal/config"
)

func CompressMiddleware(cfg *config.Config) fiber.Handler {
	level := compress.LevelBestSpeed

	switch cfg.App.CompressLevel {
	case 0:
		level = compress.LevelDisabled
	case 1:
		level = compress.LevelBestSpeed
	case 2:
		level = compress.LevelBestCompression
	default:
		level = compress.LevelBestSpeed
	}

	return compress.New(compress.Config{
		Level: level,
	})
}
