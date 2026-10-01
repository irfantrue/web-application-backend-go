package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port/mocks"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductUsecase_GetByID(t *testing.T) {
	mockRepo := new(mocks.ProductRepositoryMock)
	mockCache := new(mocks.CacheServiceMock)
	mockUoW := new(mocks.UnitOfWorkMock)

	usecase := NewProductUsecase(mockUoW, mockRepo, mockCache)

	testID, _ := uuid.NewV7()
	ctx := context.Background()

	t.Run("Success - Data found in database (Cache Miss)", func(t *testing.T) {
		mockCache.On("Get", mock.Anything, "product:"+testID.String(), mock.Anything).
			Return(errors.New("cache miss")).Once()

		mockProduct := &entity.Product{
			ID:        testID,
			Name:      "Test Product",
			Price:     decimal.NewFromInt(10000),
			CreatedAt: time.Now(),
		}
		mockRepo.On("FindByID", mock.Anything, testID).Return(mockProduct, nil).Once()

		mockCache.On("Set", mock.Anything, "product:"+testID.String(), mockProduct, mock.Anything).
			Return(nil).Once()

		res, err := usecase.FindByID(ctx, testID.String())

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "Test Product", res.Name)

		mockCache.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
