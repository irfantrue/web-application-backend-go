package port

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// ErrRefreshTokenNotFound is returned when no refresh token session matches the given hash.
var ErrRefreshTokenNotFound = errors.New("refresh token not found")

// RefreshTokenSessionRepository defines the contract for refresh token session data persistence operations.
type RefreshTokenSessionRepository interface {
	Create(ctx context.Context, session *entity.RefreshTokenSession) error
	FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshTokenSession, error)
	MarkRotated(ctx context.Context, id uuid.UUID, t time.Time) error
	RevokeByID(ctx context.Context, id uuid.UUID, t time.Time) error
	RevokeFamily(ctx context.Context, familyID uuid.UUID) error
	RevokeByUserID(ctx context.Context, userID uuid.UUID) error
}
