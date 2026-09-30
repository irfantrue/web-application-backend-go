package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/ctxval"
)

func TimezoneMiddleware(defaultTimezone string) fiber.Handler {
	defaultLoc, err := time.LoadLocation(defaultTimezone)
	if err != nil {
		defaultLoc = time.UTC
	}

	return func(c fiber.Ctx) error {
		tzHeader := c.Get("X-Timezone")

		loc := defaultLoc

		if tzHeader != "" {
			parsedLoc, err := time.LoadLocation(tzHeader)
			if err == nil {
				loc = parsedLoc
			}
		}

		newCtx := ctxval.SetTimezone(c, loc)
		c.SetContext(newCtx)

		return c.Next()
	}
}
