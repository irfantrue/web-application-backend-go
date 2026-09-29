package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
)

// ErrorHandler catches all errors returned by the handler
func ErrorHandler(c fiber.Ctx, err error) error {
	// If the error is already in the *Problem struct format (RFC 7807)
	if response.IsProblem(err) {
		prob, _ := response.From(err)
		return response.WriteProblem(c, prob)
	}

	// If the error is a built-in Fiber error (e.g., 404 Route Not Found, 405 Method Not Allowed)
	if e, ok := err.(*fiber.Error); ok {
		prob := &response.Problem{
			Title:  http.StatusText(e.Code),
			Status: e.Code,
			Detail: e.Message,
		}

		return response.WriteProblem(c, prob)
	}

	// Default catch-all for unexpected Internal Server Errors (500)
	prob := &response.Problem{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: err.Error(),
	}

	return response.WriteProblem(c, prob)
}
