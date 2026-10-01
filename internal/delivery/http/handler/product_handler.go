// Package handler provides HTTP request handlers for the application.
package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/validator"
)

// ProductHandler handles HTTP requests related to product operations.
type ProductHandler struct {
	usecase usecase.ProductUsecase
}

// NewProductHandler creates a new instance of ProductHandler with the given usecase.
func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

// Create handles the HTTP request to create a new product.
// It binds the request body, validates it, and passes it to the usecase layer.
func (h *ProductHandler) Create(c fiber.Ctx) error {
	req := new(dto.CreateProductRequest)

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

	productRes, err := h.usecase.Create(c.Context(), *req)
	if err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to Create Product",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[*dto.ProductResponse]{
		Status: http.StatusCreated,
		Data:   productRes,
	})
}

// FindAll handles the HTTP request to retrieve all products.
// It fetches the product list from the usecase layer and returns it to the client.
func (h *ProductHandler) FindAll(c fiber.Ctx) error {
	productRes, err := h.usecase.FindAll(c.Context())
	if err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to FindAll Product",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[[]dto.ProductResponse]{
		Status: http.StatusCreated,
		Data:   productRes,
	})
}

// FindByID handles the HTTP request to retrieves product.
// It fetches the product from the usecase layer and returns it to the client.
func (h *ProductHandler) FindByID(c fiber.Ctx) error {
	productRes, err := h.usecase.FindByID(c.Context(), c.Params("id"))
	if err != nil {
		return response.WriteProblem(c, &response.Problem{
			Title:  "Failed to FindByID Product",
			Status: http.StatusUnprocessableEntity,
			Detail: err.Error(),
		})
	}

	return response.Write(c, &response.Response[*dto.ProductResponse]{
		Status: http.StatusCreated,
		Data:   productRes,
	})
}
