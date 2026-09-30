package dto

import (
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

type CreateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=150"`
	Price string `json:"price" validate:"required,numeric"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}

type ProductResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	CreatedAt string `json:"created_at"`
}

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

func NewProductResponseList(products []entity.Product, loc *time.Location) []ProductResponse {
	responses := make([]ProductResponse, 0, len(products))

	for _, p := range products {
		responses = append(responses, *NewProductResponse(&p, loc))
	}

	return responses
}
