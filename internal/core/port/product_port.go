package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
}
