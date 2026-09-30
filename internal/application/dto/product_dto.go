// Package dto provides Data Transfer Objects (DTO) for request and response payloads.
package dto

import (
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// CreateProductRequest represents the payload for creating a new product.
type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=150"`
	Price string `json:"price" validate:"required,numeric"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

// ProductResponse represents the standard response format for a product.
type ProductResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
}

// NewProductResponse maps a Product entity to a ProductResponse DTO.
// It returns nil if the provided product is nil.
func NewProductResponse(product *entity.Product, loc *time.Location) *ProductResponse {
	if product == nil {
		return nil
	}
	return &ProductResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Price:     product.Price.StringFixedBank(2),
		CreatedAt: product.CreatedAt.In(loc).Format(time.RFC3339),
	}
}

// NewProductResponseList maps a slice of Product entities to a slice of ProductResponse DTOs.
func NewProductResponseList(products []entity.Product, loc *time.Location) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))

	for i := range products {
		res := NewProductResponse(&products[i], loc)
		responses = append(responses, *res)
	}

	return responses
}
