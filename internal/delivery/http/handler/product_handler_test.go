package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_Create(t *testing.T) {
	mockUsecase := new(mocks.ProductUsecaseMock)

	app := fiber.New()
	handler := NewProductHandler(mockUsecase)
	app.Post("/api/v1/products", handler.Create)

	t.Run("Product Successfully Created (Response 201)", func(t *testing.T) {
		reqBody := dto.CreateProductRequest{
			Name:  "Coffee",
			Price: "5000.00",
			Stock: 100,
		}

		expectedRes := &dto.ProductResponse{
			ID:    "some-uuid",
			Name:  "Kopi Arabica",
			Price: "50000",
		}

		mockUsecase.On("Create", mock.Anything, reqBody).Return(expectedRes, nil).Once()

		bodyBytes, _ := sonic.Marshal(reqBody)
		req := httptest.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/v1/products",
			bytes.NewReader(bodyBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{
			Timeout:       0,
			FailOnTimeout: false,
		})

		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		defer func() { _ = resp.Body.Close() }()

		var respBody map[string]any
		assert.NoError(t, sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&respBody))

		data := respBody["data"].(map[string]any)
		assert.Equal(t, "Kopi Arabica", data["name"])

		mockUsecase.AssertExpectations(t)
	})

	t.Run("Validation Failed (Response 400)", func(t *testing.T) {
		reqBody := dto.CreateProductRequest{
			Name:  "",
			Price: "50000.00",
		}

		bodyBytes, _ := sonic.Marshal(reqBody)
		req := httptest.NewRequestWithContext(
			context.Background(),
			http.MethodPost,
			"/api/v1/products",
			bytes.NewReader(bodyBytes),
		)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req, fiber.TestConfig{
			Timeout:       0,
			FailOnTimeout: false,
		})

		defer func() { _ = resp.Body.Close() }()

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
