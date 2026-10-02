package dto

import (
	"time"

	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
)

// CreateAPIKeyRequest represents the payload for creating a new api key.
type CreateAPIKeyRequest struct {
	ClientName string
	ExpiresAt  string
}

// APIKeyResponse represents the standard response format for a api key.
type APIKeyResponse struct {
	ClientName string
	KeyHash    string
	IsActive   bool
	ExpiresAt  *string
	CreatedAt  string
}

// NewAPIKeyResponse maps a APIKey entity to a APIKeyResponse DTO.
// It returns nil if the provided api key is nil.
func NewAPIKeyResponse(apiKey *entity.APIKey, loc *time.Location) *APIKeyResponse {
	if apiKey == nil {
		return nil
	}

	var expiresAt string
	if apiKey.ExpiresAt != nil {
		expiresAt = apiKey.ExpiresAt.In(loc).Format(time.RFC3339)
	}

	return &APIKeyResponse{
		ClientName: apiKey.ClientName,
		KeyHash:    apiKey.KeyHash,
		IsActive:   apiKey.IsActive,
		ExpiresAt:  &expiresAt,
		CreatedAt:  apiKey.CreatedAt.In(loc).Format(time.RFC3339),
	}
}

// NewAPIKeyResponseList maps a slice of APIKey entities to a slice of APIKeyResponse DTOs.
func NewAPIKeyResponseList(apiKeys []entity.APIKey, loc *time.Location) []APIKeyResponse {
	responses := make([]APIKeyResponse, 0, len(apiKeys))

	for i := range apiKeys {
		res := NewAPIKeyResponse(&apiKeys[i], loc)
		responses = append(responses, *res)
	}

	return responses
}
