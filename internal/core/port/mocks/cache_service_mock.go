package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

// CacheServiceMock is a mock implementation of the CacheService interface.
type CacheServiceMock struct {
	mock.Mock
}

// Set stores a value in the cache with the given TTL.
func (m *CacheServiceMock) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	args := m.Called(ctx, key, value, ttl)
	return args.Error(0)
}

// Get retrieves a value from the cache by key.
func (m *CacheServiceMock) Get(ctx context.Context, key string, dest any) error {
	args := m.Called(ctx, key, dest)
	return args.Error(0)
}

// Delete removes a key from the cache.
func (m *CacheServiceMock) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

// DeleteByPrefix removes all keys matching the given prefix.
func (m *CacheServiceMock) DeleteByPrefix(ctx context.Context, prefix string) error {
	args := m.Called(ctx, prefix)
	return args.Error(0)
}
