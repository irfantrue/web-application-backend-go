package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/validator"
)

// APIKeyHandler handles HTTP requests related to api key operations.
type APIKeyHandler struct {
	usecase usecase.APIKeyUsecase
}

// NewAPIKeyHandler creates a new instance of APIKeyHandler with the given usecase.
func NewAPIKeyHandler(u usecase.APIKeyUsecase) *APIKeyHandler {
	return &APIKeyHandler{usecase: u}
}

// Generate handles the creation of a new API key.
func (h *APIKeyHandler) Generate(c fiber.Ctx) error {
	req := new(dto.CreateAPIKeyRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Invalid Request Payload",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		})
	}

	if validationErrs := validator.ValidateStruct(req); validationErrs != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:         "Payload Validation Failed",
			Status:        http.StatusBadRequest,
			Detail:        "One or more columns do not meet the validation requirements.",
			InvalidParams: validationErrs,
		})
	}

	apiKeyRes, err := h.usecase.Generate(c.Context(), *req)
	if err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to Generate API Key",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[*dto.APIKeyResponse]{
		Status: http.StatusCreated,
		Data:   apiKeyRes,
	})
}

// FindAll retrieves all API keys for the authenticated user.
func (h *APIKeyHandler) FindAll(c fiber.Ctx) error {
	apiKeyRes, err := h.usecase.FindAll(c.Context())
	if err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to FindAll ApiKeys",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[[]dto.APIKeyResponse]{
		Status: http.StatusOK,
		Data:   apiKeyRes,
	})
}

// Activate handles the activation of an existing API key.
func (h *APIKeyHandler) Activate(c fiber.Ctx) error {
	if err := h.usecase.Activate(c.Context(), c.Params("id")); err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to Activate ApiKey",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusNoContent,
		Data:   nil,
	})
}

// Revoke handles the revocation of an existing API key.
func (h *APIKeyHandler) Revoke(c fiber.Ctx) error {
	if err := h.usecase.Revoke(c.Context(), c.Params("id")); err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to Revoke ApiKey",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusNoContent,
		Data:   nil,
	})
}
