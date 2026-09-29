package response

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
)

const ContentType = "application/json"

// Response is the standard API response envelope for all successful responses.
type Response[T any] struct {
	Status int `json:"status,omitempty"`
	Data   T   `json:"data,omitempty"`
}

// Problem is the root RFC 7807 problem detail object.
type Problem struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func (p *Problem) Error() string {
	return fmt.Sprintf("[%d] %s — %s", p.Status, p.Title, p.Detail)
}

// IsProblem reports whether err is or wraps a *Problem.
func IsProblem(err error) bool {
	var p *Problem
	return errors.As(err, &p)
}

// From extracts a *Problem from err. Returns nil if not found.
func From(err error) (*Problem, bool) {
	return errors.AsType[*Problem](err)
}

// writeJSON is an internal helper that writes a JSON response with the given status code.
//
// Returns:
//   - error: nil on success, or a wrapped error if JSON encoding/writing fails
func WriteJSON(c fiber.Ctx, status int, body any) error {
	c.Set(fiber.HeaderContentType, ContentType)
	if err := c.Status(status).JSON(body); err != nil {
		return fmt.Errorf("Failed to write JSON response")
	}
	return nil
}

// writeProblem sends a Problem (RFC 7807) response to the client.
//
// Returns:
//   - error: Result of c.JSON() operation (typically nil on success)
func WriteProblem(c fiber.Ctx, p *Problem) error {
	return WriteJSON(c, p.Status, p)
}

// Write sends a standard API response to the client.
//
// Returns:
//   - error: Result of c.JSON() operation (typically nil on success)
func Write[T any](c fiber.Ctx, r *Response[T]) error {
	c.Set(fiber.HeaderContentType, ContentType)
	res := *r
	return WriteJSON(c, res.Status, res)
}
