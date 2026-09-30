package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/database"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) port.ProductRepository {
	return &productRepository{db: db}
}

// Create implements [port.ProductRepository].
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	db := database.ExtractTx(ctx, r.db)
	dbModel := model.FromDomain(product)

	return db.WithContext(ctx).Create(dbModel).Error
}

// FindAll implements [port.ProductRepository].
func (r *productRepository) FindAll(ctx context.Context) ([]entity.Product, error) {
	panic("unimplemented")
}

// FindByID implements [port.ProductRepository].
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	panic("unimplemented")
}
