package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID        uuid.UUID
	Name      string
	Price     decimal.Decimal
	Stock     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewProduct(name string, price decimal.Decimal, stock int) (*Product, error) {
	if price.IsNegative() {
		return nil, errors.New("The price cannot be negative.")
	}

	if stock < 0 {
		return nil, errors.New("The stock cannot be negative.")
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
