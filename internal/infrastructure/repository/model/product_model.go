// Package model provides GORM database models and mapper functions to convert between domain entities and database records.
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/shopspring/decimal"
)

// ProductModel represents the database table schema for products in GORM.
type ProductModel struct {
	ID          uuid.UUID       `gorm:"primaryKey;column:id"`
	Name        string          `gorm:"type:varchar(100);column:name"`
	Description string          `gorm:"type:text"`
	Price       decimal.Decimal `gorm:"type:decimal(18,2);column:price"`
	Stock       int             `gorm:"not null;default:0;column:stock"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

// TableName specifies the custom table name for ProductModel in GORM.
func (ProductModel) TableName() string {
	return "products"
}

// FromDomain converts a core domain Product entity into a database ProductModel.
func FromDomain(p *entity.Product) *ProductModel {
	return &ProductModel{
		ID:        p.ID,
		Name:      p.Name,
		Price:     p.Price,
		Stock:     p.Stock,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToDomain converts a database ProductModel record into a core domain Product entity.
func (m *ProductModel) ToDomain() *entity.Product {
	return &entity.Product{
		ID:        m.ID,
		Name:      m.Name,
		Price:     m.Price,
		Stock:     m.Stock,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
