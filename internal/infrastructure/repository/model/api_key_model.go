package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// APIKeyModel represents the database table schema for api keys in GORM.
type APIKeyModel struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ClientName string     `gorm:"column:client_name;type:varchar(255);not null"`
	KeyHash    string     `gorm:"column:key_hash;type:varchar(255);uniqueIndex;not null"`
	IsActive   bool       `gorm:"column:is_active;type:bool;not null;default:true"`
	ExpiresAt  *time.Time `gorm:"column:expires_at" `
	CreatedAt  time.Time  `gorm:"column:created_at"`
	UpdatedAt  time.Time  `gorm:"column:updated_at"`
}

// TableName specifies the custom table name for APIKeyModel in GORM.
func (APIKeyModel) TableName() string {
	return "api_keys"
}

// NewAPIKeyModel converts a core domain APIKey entity into a database APIKeyModel.
func NewAPIKeyModel(p *entity.APIKey) *APIKeyModel {
	return &APIKeyModel{
		ID:         p.ID,
		ClientName: p.ClientName,
		KeyHash:    p.KeyHash,
		IsActive:   p.IsActive,
		ExpiresAt:  p.ExpiresAt,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

// ToDomain converts a database APIKeyModel record into a core domain APIKey entity.
func (m *APIKeyModel) ToDomain() *entity.APIKey {
	return &entity.APIKey{
		ID:         m.ID,
		ClientName: m.ClientName,
		KeyHash:    m.KeyHash,
		IsActive:   m.IsActive,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
