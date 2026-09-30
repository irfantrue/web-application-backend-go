package database

import (
	"context"

	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"gorm.io/gorm"
)

type txKey struct{}

type GormUoW struct {
	db *gorm.DB
}

func NewGormUnitOfWork(db *gorm.DB) port.UnitOfWork {
	return &GormUoW{db: db}
}

func (u *GormUoW) Execute(ctx context.Context, fn func(txCtx context.Context) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func ExtractTx(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
