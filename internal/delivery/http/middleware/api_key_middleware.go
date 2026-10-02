package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/application/usecase"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
)

// APIKeyMiddleware returns a Fiber handler that validates the Api Key
// from the Authorization header.
func APIKeyMiddleware(apiKeyUsecase usecase.APIKeyUsecase) fiber.Handler {
	return func(c fiber.Ctx) error {
		clientKey := c.Get("X-API-KEY")

		if clientKey == "" {
			return response.WriteProblem(c, &response.Problem{
				Title:  "Missing API Key",
				Status: http.StatusUnauthorized,
				Detail: "The X-API-Key header was not found.",
			})
		}

		if err := apiKeyUsecase.Validate(c.Context(), clientKey); err != nil {
			return response.WriteProblem(c, &response.Problem{
				Title:  "Invalid API Key",
				Status: http.StatusUnauthorized,
				Detail: err.Error(),
			})
		}

		return c.Next()
	}
}
