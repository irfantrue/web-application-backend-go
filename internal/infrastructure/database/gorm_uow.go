// Package database provides infrastructure implementations for database operations.
package database

import (
	"context"

	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"gorm.io/gorm"
)

// txKey is an unexported custom type used as a key for storing the GORM transaction in the context.
// Using an unexported type prevents collisions with other context values.
type txKey struct{}

// GormUoW implements the port.UnitOfWork interface using GORM.
type GormUoW struct {
	db *gorm.DB
}

// NewGormUnitOfWork creates a new instance of GormUoW.
func NewGormUnitOfWork(db *gorm.DB) port.UnitOfWork {
	return &GormUoW{db: db}
}

// Execute runs the provided function within a GORM database transaction.
// It injects the transaction object into the context so repositories can use it.
func (u *GormUoW) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// ExtractTx retrieves the GORM transaction from the context.
// If no transaction is found in the context, it gracefully falls back to the defaultDB.
func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
