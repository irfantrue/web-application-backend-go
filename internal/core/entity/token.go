package entity

import (
	"time"
)

type TokenPayload struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	TokenVersion int    `json:"token_version"`
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
