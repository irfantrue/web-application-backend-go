package http

import (
	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/handler"
)

type SetupRoutes struct {
	App           *fiber.App
	HealthHandler *handler.HealthHandler
}

func (r *SetupRoutes) Build() {
	r.App.Get("/health", r.HealthHandler.Get)
}
