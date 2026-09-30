package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/shopspring/decimal"
)

type ProductModel struct {
	ID          uuid.UUID       `gorm:"primaryKey;column:id"`
	Name        string          `gorm:"type:varchar(100);column:name"`
	Description string          `gorm:"type:text"`
	Price       decimal.Decimal `gorm:"type:decimal(18,2);column:price"`
	Stock       int             `gorm:"not null;default:0;column:stock"`
	CreatedAt   time.Time       `gorm:"column:created_at"`
	UpdatedAt   time.Time       `gorm:"column:updated_at"`
}

func (ProductModel) TableName() string {
	return "products"
}

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
