package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/database"
	"github.com/irfantrue/web-application-backend-go/internal/infrastructure/repository/model"
	"gorm.io/gorm"
)

var _ port.RefreshTokenSessionRepository = (*rtsRepository)(nil)

type rtsRepository struct {
	db *gorm.DB
}

// NewRefreshTokenSessionRepository creates a new RefreshTokenSessionRepository.
func NewRefreshTokenSessionRepository(db *gorm.DB) port.RefreshTokenSessionRepository {
	return &rtsRepository{db: db}
}

// Create inserts a new refresh token session.
func (r *rtsRepository) Create(ctx context.Context, session *entity.RefreshTokenSession) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).Create(model.NewRefreshTokenSessionModel(session)).Error
}

// FindByTokenHash retrieves a session by its token hash.
func (r *rtsRepository) FindByTokenHash(ctx context.Context, hash string) (*entity.RefreshTokenSession, error) {
	db := database.ExtractTx(ctx, r.db)

	var model model.RefreshTokenSessionModel
	if err := db.WithContext(ctx).
		Where("token_hash = ?", hash).
		First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrRefreshTokenNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

// MarkRotated sets the rotated_at timestamp on a session.
func (r *rtsRepository) MarkRotated(ctx context.Context, id uuid.UUID, t time.Time) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.RefreshTokenSessionModel{}).
		Where("id = ?", id).
		UpdateColumn("rotated_at", t).
		Error
}

// RevokeByID revokes a single session.
func (r *rtsRepository) RevokeByID(ctx context.Context, id uuid.UUID, t time.Time) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.RefreshTokenSessionModel{}).
		Where("id = ? AND revoked_at IS NULL", id).
		UpdateColumn("revoked_at", t).
		Error
}

// RevokeFamily revokes all sessions in a token family (reuse detection).
func (r *rtsRepository) RevokeFamily(ctx context.Context, familyID uuid.UUID) error {
	db := database.ExtractTx(ctx, r.db)

	now := time.Now()
	return db.WithContext(ctx).
		Model(&model.RefreshTokenSessionModel{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		UpdateColumn("revoked_at", now).
		Error
}

// RevokeByUserID revokes all active sessions for a user.
func (r *rtsRepository) RevokeByUserID(ctx context.Context, userID uuid.UUID) error {
	db := database.ExtractTx(ctx, r.db)

	now := time.Now()
	return db.WithContext(ctx).
		Model(&model.RefreshTokenSessionModel{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		UpdateColumn("revoked_at", now).
		Error
}
