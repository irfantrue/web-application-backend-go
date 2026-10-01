package mocks

import (
	"context"

	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/stretchr/testify/mock"
)

// ProductUsecaseMock is a mock implementation of the ProductUsecase interface.
type ProductUsecaseMock struct {
	mock.Mock
}

// Create creates a new product.
func (m *ProductUsecaseMock) Create(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	args := m.Called(ctx, req)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

// FindByID returns a product by its ID.
func (m *ProductUsecaseMock) FindByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProductResponse), args.Error(1)
}

// FindAll returns all products.
func (m *ProductUsecaseMock) FindAll(ctx context.Context) ([]dto.ProductResponse, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ProductResponse), args.Error(1)
}
