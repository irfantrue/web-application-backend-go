package port

import (
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// TokenService handles the generation and validation of authentication tokens.
type TokenService interface {
	// GenerateAccessToken creates a new token for the specified user ID.
	GenerateAccessToken(payload *entity.TokenPayload) (*entity.TokenResponse, error)

	// ValidateToken checks the validity of a token and returns the associated user ID if valid.
	ValidateToken(token string) (*entity.TokenPayload, error)
}
