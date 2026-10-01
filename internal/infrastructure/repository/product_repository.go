// Package repository provides GORM-based implementations of the core application ports.
package repository

import (
	"context"
	"errors"

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

// NewProductRepository creates a new instance of port.ProductRepository backed by GORM.
func NewProductRepository(db *gorm.DB) port.ProductRepository {
	return &productRepository{db: db}
}

// Create implements [port.ProductRepository].
func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	db := database.ExtractTx(ctx, r.db)
	dbModel := model.NewProductModel(product)

	return db.WithContext(ctx).Create(dbModel).Error
}

// FindAll implements [port.ProductRepository].
func (r *productRepository) FindAll(ctx context.Context) ([]entity.Product, error) {
	db := database.ExtractTx(ctx, r.db)

	dbModels := []model.ProductModel{}
	if err := db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}

	products := make([]entity.Product, 0, len(dbModels))
	for i := range dbModels {
		domainEntity := dbModels[i].ToDomain()
		if domainEntity != nil {
			products = append(products, *domainEntity)
		}
	}

	return products, nil
}

// FindByID implements [port.ProductRepository].
func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	db := database.ExtractTx(ctx, r.db)

	dbModel := model.ProductModel{}
	if err := db.WithContext(ctx).Where("id = ?", id).First(&dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dbModel.ToDomain(), nil
}
