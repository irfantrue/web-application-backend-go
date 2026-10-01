// Package entity contains the core business domain models and logic.
package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Product represents an item available for sale in the system.
type Product struct {
	ID        uuid.UUID
	Name      string
	Price     decimal.Decimal
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProduct creates a new Product instance with validation and default values.
// It returns an error if the price or stock is negative.
func NewProduct(name string, price decimal.Decimal, stock int) (*Product, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	if price.IsNegative() {
		return nil, errors.New("price cannot be negative")
	}

	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}

	now := time.Now()

	return &Product{
		ID:        uuid.New(),
		Name:      name,
		Price:     price,
		Stock:     stock,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
