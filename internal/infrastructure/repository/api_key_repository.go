package repository

import (
	"context"
	"errors"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/database"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

type apiKeyRepository struct {
	db *gorm.DB
}

// NewAPIKeyRepository creates a new instance of port.APIKeyRepository backed by GORM.
func NewAPIKeyRepository(db *gorm.DB) port.APIKeyRepository {
	return &apiKeyRepository{db: db}
}

// Create implements [port.APIKeyRepository].
func (r *apiKeyRepository) Create(ctx context.Context, apiKey *entity.APIKey) error {
	db := database.ExtractTx(ctx, r.db)
	dbModel := model.NewAPIKeyModel(apiKey)

	return db.WithContext(ctx).Create(dbModel).Error
}

// FindByKeyString implements [port.APIKeyRepository].
func (r *apiKeyRepository) FindByKeyString(ctx context.Context, apiKey string) (*entity.APIKey, error) {
	db := database.ExtractTx(ctx, r.db)

	var dbModel model.APIKeyModel
	if err := db.WithContext(ctx).Where("api_key = ?", apiKey).First(&dbModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return dbModel.ToDomain(), nil
}

// FindAll implements [port.APIKeyRepository].
func (r *apiKeyRepository) FindAll(ctx context.Context) ([]entity.APIKey, error) {
	db := database.ExtractTx(ctx, r.db)

	var dbModels []model.APIKeyModel
	if err := db.WithContext(ctx).Find(&dbModels).Error; err != nil {
		return nil, err
	}

	apiKeys := make([]entity.APIKey, 0, len(dbModels))
	for i := range dbModels {
		domainEntity := dbModels[i].ToDomain()
		if domainEntity != nil {
			apiKeys = append(apiKeys, *domainEntity)
		}
	}

	return apiKeys, nil
}

func (r *apiKeyRepository) setActive(ctx context.Context, keyString string, isActive bool) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.APIKeyModel{}).
		Where("key_string = ?", keyString).
		UpdateColumn("is_active", isActive).Error
}

// Activate implements [port.APIKeyRepository].
func (r *apiKeyRepository) Activate(ctx context.Context, keyString string) error {
	return r.setActive(ctx, keyString, true)
}

// Revoke implements [port.APIKeyRepository].
func (r *apiKeyRepository) Revoke(ctx context.Context, keyString string) error {
	return r.setActive(ctx, keyString, false)
}
