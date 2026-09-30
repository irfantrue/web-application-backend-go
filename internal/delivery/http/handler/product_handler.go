package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/validator"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

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
