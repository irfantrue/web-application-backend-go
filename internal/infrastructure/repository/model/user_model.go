package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// UserModel represents the database table schema for users in GORM.
type UserModel struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	Email               string     `gorm:"type:varchar(255);uniqueIndex;not null;column:email"`
	Password            string     `gorm:"type:varchar(255);not null;column:password"`
	TokenVersion        int        `gorm:"not null;default:0;column:token_version"`
	ResetTokenHash      string     `gorm:"type:varchar(255);column:reset_token_hash"`
	ResetTokenExpiresAt *time.Time `gorm:"column:reset_token_expires_at"`
	CreatedAt           time.Time  `gorm:"column:created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at"`
	LastLoginAt         *time.Time `gorm:"column:last_login_at"`
}

// TableName specifies the custom table name for UserModel in GORM.
func (UserModel) TableName() string {
	return "users"
}

// NewUserModel converts a core domain User entity into a database UserModel.
func NewUserModel(p *entity.User) *UserModel {
	return &UserModel{
		ID:                  p.ID,
		Email:               p.Email,
		Password:            p.Password,
		TokenVersion:        p.TokenVersion,
		ResetTokenHash:      p.ResetTokenHash,
		ResetTokenExpiresAt: p.ResetTokenExpiresAt,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
		LastLoginAt:         p.LastLoginAt,
	}
}

// ToDomain converts a database UserModel record into a core domain User entity.
func (m *UserModel) ToDomain() *entity.User {
	return &entity.User{
		ID:                  m.ID,
		Email:               m.Email,
		Password:            m.Password,
		TokenVersion:        m.TokenVersion,
		ResetTokenHash:      m.ResetTokenHash,
		ResetTokenExpiresAt: m.ResetTokenExpiresAt,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		LastLoginAt:         m.LastLoginAt,
	}
}
