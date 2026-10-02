package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/irfantrue/web-application-backend-go/internal/pkg/ctxval"
)

var (
	// ErrAPIKeyNotFound is returned when the requested API key does not exist.
	ErrAPIKeyNotFound = errors.New("api key not found")

	// ErrAPIKeyRevoked is returned when the API key has been revoked.
	ErrAPIKeyRevoked = errors.New("api key has been revoked")

	// ErrAPIKeyExpired is returned when the API key has passed its expiration time.
	ErrAPIKeyExpired = errors.New("api key has expired")
)

// APIKeyUsecase defines the business logic interface for api key operations.
type APIKeyUsecase interface {
	Generate(ctx context.Context, req dto.CreateAPIKeyRequest) (*dto.APIKeyResponse, error)
	FindAll(ctx context.Context) ([]dto.APIKeyResponse, error)
	Activate(ctx context.Context, keyHash string) error
	Revoke(ctx context.Context, keyHash string) error
	Validate(ctx context.Context, keyHash string) error
}

type apiKeyUsecaseImpl struct {
	uow    port.UnitOfWork
	repo   port.APIKeyRepository
	cache  port.CacheService
	crypto port.CryptographyService
}

// NewAPIKeyUsecase creates a new instance of APIKeyUsecase.
func NewAPIKeyUsecase(
	uow port.UnitOfWork,
	repo port.APIKeyRepository,
	cache port.CacheService,
	crypto port.CryptographyService,
) APIKeyUsecase {
	return &apiKeyUsecaseImpl{
		uow:    uow,
		repo:   repo,
		cache:  cache,
		crypto: crypto,
	}
}

// Generate implements [APIKeyUsecase].
func (u *apiKeyUsecaseImpl) Generate(ctx context.Context, req dto.CreateAPIKeyRequest) (*dto.APIKeyResponse, error) {
	randomStr, err := u.crypto.GenerateSecureString(32)
	if err != nil {
		return nil, err
	}

	rawKey := "pk_live_" + randomStr
	hashedKey := u.crypto.HashString(rawKey)

	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("parse expires_at: %w", err)
		}
		if t.Before(time.Now()) {
			return nil, errors.New("expires_at must be in the future")
		}
		expiresAt = &t
	}

	apiKey := &entity.APIKey{
		ID:         uuid.New(),
		KeyHash:    hashedKey,
		ClientName: req.ClientName,
		IsActive:   true,
		ExpiresAt:  expiresAt,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := u.repo.Create(ctx, apiKey); err != nil {
		return nil, fmt.Errorf("create api key: %w", err)
	}

	userLoc := ctxval.GetTimezone(ctx)
	return dto.NewAPIKeyResponse(apiKey, userLoc), nil
}

// FindAll implements [APIKeyUsecase].
func (u *apiKeyUsecaseImpl) FindAll(ctx context.Context) ([]dto.APIKeyResponse, error) {
	apiKeys, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve api keys: %w", err)
	}

	userLoc := ctxval.GetTimezone(ctx)
	return dto.NewAPIKeyResponseList(apiKeys, userLoc), nil
}

// Activate implements [APIKeyUsecase].
func (u *apiKeyUsecaseImpl) Activate(ctx context.Context, keyHash string) error {
	if err := u.repo.Activate(ctx, keyHash); err != nil {
		return fmt.Errorf("activate api key: %w", err)
	}

	_ = u.cache.Set(ctx, cacheKey(keyHash), true, 1*time.Hour)
	return nil
}

// Revoke implements [APIKeyUsecase].
func (u *apiKeyUsecaseImpl) Revoke(ctx context.Context, keyHash string) error {
	if err := u.repo.Revoke(ctx, keyHash); err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}

	_ = u.cache.Delete(ctx, cacheKey(keyHash))
	return nil
}

// Validate implements [APIKeyUsecase].
func (u *apiKeyUsecaseImpl) Validate(ctx context.Context, keyHash string) error {
	ck := cacheKey(keyHash)

	var cached bool
	if err := u.cache.Get(ctx, ck, &cached); err != nil {
		if cached {
			return nil
		}
		return ErrAPIKeyRevoked
	}

	apiKey, err := u.repo.FindByKeyString(ctx, keyHash)
	if err != nil {
		return ErrAPIKeyNotFound
	}

	if !apiKey.IsActive {
		_ = u.cache.Set(ctx, ck, false, 5*time.Minute)
		return ErrAPIKeyRevoked
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		_ = u.cache.Set(ctx, ck, false, 5*time.Minute)
		return ErrAPIKeyExpired
	}

	_ = u.cache.Set(ctx, ck, true, 1*time.Hour)

	return nil
}

func cacheKey(keyHash string) string {
	return "apikey_valid:" + keyHash
}
