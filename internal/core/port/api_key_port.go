package port

import (
	"context"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// APIKeyRepository defines the contract for user data persistence operations.
type APIKeyRepository interface {
	Create(ctx context.Context, apiKey *entity.APIKey) error
	FindByKeyString(ctx context.Context, keyHash string) (*entity.APIKey, error)
	FindAll(ctx context.Context) ([]entity.APIKey, error)
	Activate(ctx context.Context, keyHash string) error
	Revoke(ctx context.Context, keyHash string) error
}
