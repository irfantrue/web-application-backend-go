package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
)

func RecoverMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] %v\n%s\n", r, debug.Stack())

				_ = response.WriteProblem(c, &response.Problem{
					Title:  "Internal Server Error",
					Status: http.StatusInternalServerError,
					Detail: "An internal server error has occurred. Our team has logged this issue.",
				})
			}
		}()

		return c.Next()
	}
}
