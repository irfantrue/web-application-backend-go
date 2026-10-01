// Package mocks provides mock implementations of repository interfaces.
package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/stretchr/testify/mock"
)

// ProductRepositoryMock is a mock implementation of the ProductRepository interface.
type ProductRepositoryMock struct {
	mock.Mock
}

// Create stores a new product.
func (m *ProductRepositoryMock) Create(ctx context.Context, product *entity.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

// FindAll returns all products.
func (m *ProductRepositoryMock) FindAll(ctx context.Context) ([]entity.Product, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Product), args.Error(1)
}

// FindByID returns a product by its ID.
func (m *ProductRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error) {
	args := m.Called(ctx, id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*entity.Product), args.Error(1)
}
