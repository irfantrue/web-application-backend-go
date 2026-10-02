package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// User represents an user available in the system.
type User struct {
	ID                  uuid.UUID
	Email               string
	Password            string
	TokenVersion        int
	ResetTokenHash      string
	ResetTokenExpiresAt *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	LastLoginAt         *time.Time
}

// NewUser creates a new User instance with validation and default values.
// It returns an error if the price or stock is negative.
func NewUser(email, password string) (*User, error) {
	now := time.Now()

	userID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed generate user id: %w", err)
	}

	return &User{
		ID:           userID,
		Email:        email,
		Password:     password,
		TokenVersion: 0,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
