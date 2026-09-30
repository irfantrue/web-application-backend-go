// Package response provides standard HTTP response formats, including success envelopes and RFC 7807 problem details.
package response

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// ContentType is the default media type for JSON responses.
const ContentType = "application/json"

// Response is the standard API response envelope for all successful responses.
type Response[T any] struct {
	Status int `json:"status,omitempty"`
	Data   T   `json:"data,omitempty"`
}

// InvalidParam merepresentasikan detail error validasi per kolom
type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

// Problem is the root RFC 7807 problem detail object.
type Problem struct {
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail,omitempty"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

// Error implements the error interface for the Problem struct.
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

// WriteJSON is an internal helper that writes a JSON response with the given status code.
func WriteJSON(c fiber.Ctx, status int, body any) error {
	c.Set(fiber.HeaderContentType, ContentType)
	if err := c.Status(status).JSON(body); err != nil {
		return fmt.Errorf("failed to write JSON response")
	}
	return nil
}

// WriteProblem sends a Problem (RFC 7807) response to the client.
func WriteProblem(c fiber.Ctx, p *Problem) error {
	return WriteJSON(c, p.Status, p)
}

// Write sends a standard API response to the client.
func Write[T any](c fiber.Ctx, r *Response[T]) error {
	c.Set(fiber.HeaderContentType, ContentType)
	res := *r
	return WriteJSON(c, res.Status, res)
}

// UnauthorizedResponse returns a standard RFC 7807 problem detail response for unauthorized access (401).
func UnauthorizedResponse(c fiber.Ctx, detail string) error {
	return WriteProblem(c, &Problem{
		Title:  "Unauthorized Access",
		Status: http.StatusUnauthorized,
		Detail: detail,
	})
}
