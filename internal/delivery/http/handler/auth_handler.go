package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/validator"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: u}
}

func (h *AuthHandler) Register(c fiber.Ctx) error {
	req := new(dto.RegisterRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	if validationErrs := validator.ValidateStruct(req); validationErrs != nil {
		return response.BadRequestResponse(
			c,
			"Payload Validation Failed",
			"One or more columns do not meet the validation requirements.",
			validationErrs,
		)
	}

	if err := h.usecase.Register(c.Context(), *req); err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to Register", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusCreated,
		Data:   nil,
	})
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	req := new(dto.LoginRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	loginRes, err := h.usecase.Login(c.Context(), *req)
	if err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to Login", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusOK,
		Data:   loginRes,
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	req := new(dto.LogoutRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	if err := h.usecase.Logout(c.Context(), req.RefreshToken); err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to Logout", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusNoContent,
		Data:   nil,
	})
}

func (h *AuthHandler) Refresh(c fiber.Ctx) error {
	req := new(dto.RefreshRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	refreshRes, err := h.usecase.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to Refresh", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusOK,
		Data:   refreshRes,
	})
}

func (h *AuthHandler) ForgotPassword(c fiber.Ctx) error {
	req := new(dto.ForgotPasswordRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	if err := h.usecase.ForgotPassword(c.Context(), req.Email); err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to ForgotPassword", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusOK,
		Data:   nil,
	})
}

func (h *AuthHandler) ResetPassword(c fiber.Ctx) error {
	req := new(dto.ResetPasswordRequest)

	if err := c.Bind().Body(req); err != nil {
		return response.BadRequestResponse(c, "Invalid Request Payload", err.Error(), nil)
	}

	if err := h.usecase.ResetPassword(c.Context(), *req); err != nil {
		return response.UnprocessableEntityResponse(c, "Failed to ResetPassword", err.Error())
	}

	return response.Write(c, &response.Response[any]{
		Status: http.StatusNoContent,
		Data:   nil,
	})
}
