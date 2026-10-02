package dto

import (
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// RegisterRequest represents the payload for creating a new user.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=8,max=32"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=32"`
}

// AuthResponse represents the standard response format for a user.
type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type MeResponse struct{}

// NewAuthResponse maps a User entity to a AuthResponse DTO.
// It returns nil if the provided user is nil.
func NewAuthResponse(accessToken, refreshToken string, expiresIn time.Duration) *AuthResponse {
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(expiresIn.Seconds()),
	}
}

func NewMeResponse(user *entity.User) *MeResponse {
	if user == nil {
		return nil
	}

	return &MeResponse{}
}
