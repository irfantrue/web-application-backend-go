package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// RefreshTokenSessionModel represents the database table schema for users in GORM.
type RefreshTokenSessionModel struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	FamilyID  uuid.UUID  `gorm:"type:uuid;not null;index;column:family_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index;column:user_id"`
	TokenHash string     `gorm:"type:varchar(255);uniqueIndex;not null;column:token_hash"`
	DeviceID  string     `gorm:"type:varchar(255);column:device_id"`
	IP        string     `gorm:"type:varchar(45);column:ip"`
	UserAgent string     `gorm:"type:varchar(255);column:user_agent"`
	RotatedAt *time.Time `gorm:"column:rotated_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`
	ExpiresAt time.Time  `gorm:"not null;column:expires_at"`
	CreatedAt time.Time  `gorm:"column:created_at"`
}

// TableName specifies the custom table name for RefreshTokenSessionModel in GORM.
func (RefreshTokenSessionModel) TableName() string {
	return "refresh_token_session"
}

// NewRefreshTokenSessionModel converts a domain entity into a database model.
func NewRefreshTokenSessionModel(p *entity.RefreshTokenSession) *RefreshTokenSessionModel {
	return &RefreshTokenSessionModel{
		ID:        p.ID,
		FamilyID:  p.FamilyID,
		UserID:    p.UserID,
		TokenHash: p.TokenHash,
		DeviceID:  p.DeviceID,
		IP:        p.IP,
		UserAgent: p.UserAgent,
		RotatedAt: p.RotatedAt,
		RevokedAt: p.RevokedAt,
		ExpiresAt: p.ExpiresAt,
		CreatedAt: p.CreatedAt,
	}
}

// ToDomain converts a database model into a domain entity.
func (m *RefreshTokenSessionModel) ToDomain() *entity.RefreshTokenSession {
	return &entity.RefreshTokenSession{
		ID:        m.ID,
		FamilyID:  m.FamilyID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		DeviceID:  m.DeviceID,
		IP:        m.IP,
		UserAgent: m.UserAgent,
		RotatedAt: m.RotatedAt,
		RevokedAt: m.RevokedAt,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}
