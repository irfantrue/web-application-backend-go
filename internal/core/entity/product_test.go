package entity

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNewProduct(t *testing.T) {
	tests := []struct {
		name          string
		productName   string
		productPrice  decimal.Decimal
		productStock  int
		expectedError string
	}{
		{
			name:          "Sukses membuat produk",
			productName:   "Laptop Asus",
			productPrice:  decimal.NewFromInt(15000000),
			productStock:  10,
			expectedError: "",
		},
		{
			name:          "Gagal - Nama kosong",
			productName:   "",
			productPrice:  decimal.NewFromInt(10000),
			productStock:  5,
			expectedError: "name cannot be empty",
		},
		{
			name:          "Gagal - Harga negatif",
			productName:   "Mouse",
			productPrice:  decimal.NewFromInt(-50000),
			productStock:  5,
			expectedError: "price cannot be negative",
		},
		{
			name:          "Gagal - Stok negatif",
			productName:   "Keyboard",
			productPrice:  decimal.NewFromInt(100000),
			productStock:  -2,
			expectedError: "stock cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.productPrice, tt.productStock)

			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, product)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tt.productName, product.Name)
				assert.NotNil(t, product.ID)
			}
		})
	}
}
