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

type ProductUsecase interface {
	Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error)
}

type productUsecase struct {
	uow         port.UnitOfWork
	productRepo port.ProductRepository
}

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
		return nil, fmt.Errorf("Invalid price format: %w", err)
	}

	product, err := entity.NewProduct(req.Name, priceDec, req.Stock)
	if err != nil {
		return nil, fmt.Errorf("Failed to create product entity: %w", err)
	}

	err = u.uow.Execute(ctx, func(txCtx context.Context) error {
		return u.productRepo.Create(txCtx, product)
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to save product: %w", err)
	}

	userLoc := ctxval.GetTimezone(ctx)

	return dto.NewProductResponse(product, userLoc), nil
}
