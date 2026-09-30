// Package usecase provides the business logic layer for the application.
package usecase

import (
	"context"
	"fmt"

	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/ctxval"
	"github.com/shopspring/decimal"
)

// ProductUsecase defines the business logic interface for product operations.
type ProductUsecase interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
	FindAll(ctx context.Context) ([]dto.ProductResponse, error)
}

type productUsecase struct {
	uow         port.UnitOfWork
	productRepo port.ProductRepository
}

// NewProductUsecase creates a new instance of ProductUsecase.
func NewProductUsecase(uow port.UnitOfWork, productRepo port.ProductRepository) ProductUsecase {
	return &productUsecase{
		uow:         uow,
		productRepo: productRepo,
	}
}

// Create implements [ProductUsecase].
func (u *productUsecase) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	priceDec, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	product, err := entity.NewProduct(req.Name, priceDec, req.Stock)
	if err != nil {
		return nil, fmt.Errorf("failed to create product entity: %w", err)
	}

	err = u.uow.Execute(ctx, func(txCtx context.Context) error {
		return u.productRepo.Create(txCtx, product)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save product: %w", err)
	}

	userLoc := ctxval.GetTimezone(ctx)

	return dto.NewProductResponse(product, userLoc), nil
}

func (u *productUsecase) FindAll(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}

	userLoc := ctxval.GetTimezone(ctx)

	return dto.NewProductResponseList(products, userLoc), nil
}
