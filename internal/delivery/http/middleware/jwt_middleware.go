package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/delivery/http/response"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/ctxval"
)

// JwtMiddleware returns a Fiber handler that validates the JWT Bearer token
// from the Authorization header and injects the user ID into the request context.
func JwtMiddleware(tokenService port.TokenService) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.UnauthorizedResponse(c, "Authorization header (Bearer) not found")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.UnauthorizedResponse(c, "The token format must be 'Bearer <token>'")
		}
		tokenString := parts[1]

		tokenPayload, err := tokenService.ValidateToken(tokenString)
		if err != nil {
			return response.UnauthorizedResponse(c, err.Error())
		}

		newCtx := ctxval.SetUserID(c, tokenPayload.UserID)
		c.SetContext(newCtx)

		return c.Next()
	}
}
