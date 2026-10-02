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

var _ port.UserRepository = (*userRepository)(nil)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of port.UserRepository backed by GORM.
func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{db: db}
}

// Create implements [port.UserRepository].
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).Create(model.NewUserModel(user)).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	db := database.ExtractTx(ctx, r.db)

	var model model.UserModel
	if err := db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	db := database.ExtractTx(ctx, r.db)

	var model model.UserModel
	if err := db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *userRepository) FindByResetToken(ctx context.Context, hash string) (*entity.User, error) {
	db := database.ExtractTx(ctx, r.db)

	var model model.UserModel

	err := db.WithContext(ctx).Where("reset_token_hash = ?", hash).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, port.ErrUserNotFound
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, id uuid.UUID, hashed string) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("password", hashed).
		Error
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID, t time.Time) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", id).
		UpdateColumn("last_login_at", t).
		Error
}

func (r *userRepository) IncrementTokenVersion(ctx context.Context, id uuid.UUID) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", id).
		UpdateColumn("token_version", gorm.Expr("token_version + 1")).
		Error
}

func (r *userRepository) SetResetToken(ctx context.Context, id uuid.UUID, hash string, expiresAt time.Time) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"reset_token_hash":       hash,
			"reset_token_expires_at": expiresAt,
		}).Error
}

func (r *userRepository) ClearResetToken(ctx context.Context, id uuid.UUID) error {
	db := database.ExtractTx(ctx, r.db)

	return db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"reset_token_hash":       "",
			"reset_token_expires_at": nil,
		}).Error
}
