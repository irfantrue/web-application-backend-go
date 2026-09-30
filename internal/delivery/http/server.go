package http

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/middleware"
)

// NewFiber initializes and configures a new Fiber application instance.
// It uses sonic for high-performance JSON encoding/decoding and sets up a custom error handler.
func NewFiber(cfg *config.Config) *fiber.App {
	return fiber.New(fiber.Config{
		AppName:      cfg.App.AppName,
		ServerHeader: cfg.App.ServerHeader,
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: middleware.ErrorHandler,
	})
}
