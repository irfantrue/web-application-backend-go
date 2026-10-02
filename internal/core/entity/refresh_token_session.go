package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenSession represents an refresh token session available in the system.
type RefreshTokenSession struct {
	ID        uuid.UUID
	FamilyID  uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	DeviceID  string
	IP        string
	UserAgent string
	RotatedAt *time.Time
	RevokedAt *time.Time
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewRefreshTokenSession creates a new Refresh Token Session instance with validation and default values.
func NewRefreshTokenSession(userID uuid.UUID, tokenHash string, expiresAt time.Time) (*RefreshTokenSession, error) {
	now := time.Now()

	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed generate refresh token session id: %w", err)
	}

	familyID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed generate refresh token session family id: %w", err)
	}

	return &RefreshTokenSession{
		ID:        id,
		FamilyID:  familyID,
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	}, nil
}
