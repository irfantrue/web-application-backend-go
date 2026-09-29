package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
)

type HealthHandler struct {
	usecase usecase.HealthUsecase
}

func NewHealthHandler(u usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{usecase: u}
}

func (h *HealthHandler) Get(c fiber.Ctx) error {
	health, err := h.usecase.Check(c.Context())
	if err != nil {
		return response.WriteProblem(c, &response.Problem{})
	}

	return response.Write(c, &response.Response[*dto.HealthResponse]{
		Status: http.StatusOK,
		Data:   health,
	})
}
