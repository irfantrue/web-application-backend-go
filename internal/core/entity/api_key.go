package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// APIKey represents an api key available in the system.
type APIKey struct {
	ID         uuid.UUID
	ClientName string
	KeyHash    string
	IsActive   bool
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// NewAPIKey creates a new api key instance with validation and default values.
func NewAPIKey(clientName, keyHash string, isActive bool, expiresAt *time.Time) (*APIKey, error) {
	if clientName == "" {
		return nil, errors.New("client name cannot be empty")
	}

	if keyHash == "" {
		return nil, errors.New("key string cannot be empty")
	}

	now := time.Now()

	return &APIKey{
		ID:         uuid.New(),
		ClientName: clientName,
		KeyHash:    keyHash,
		IsActive:   isActive,
		ExpiresAt:  expiresAt,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}
