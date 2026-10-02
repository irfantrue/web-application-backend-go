package port

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// ErrUserNotFound is returned when user not found.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the contract for user data persistence operations.
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByResetToken(ctx context.Context, hash string) (*entity.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, hashed string) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error
	IncrementTokenVersion(ctx context.Context, id uuid.UUID) error
	SetResetToken(ctx context.Context, id uuid.UUID, hash string, expiresAt time.Time) error
	ClearResetToken(ctx context.Context, id uuid.UUID) error
}
