package cache

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"github.com/redis/go-redis/v9"
)

type redisCacheImpl struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) port.CacheService {
	return &redisCacheImpl{client: client}
}

// Set implements [port.CacheService].
func (r *redisCacheImpl) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	bytes, err := sonic.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, bytes, ttl).Err()
}

// Get implements [port.CacheService].
func (r *redisCacheImpl) Get(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return sonic.Unmarshal([]byte(val), dest)
}

// Delete implements [port.CacheService].
func (r *redisCacheImpl) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// DeleteByPrefix implements [port.CacheService].
func (r *redisCacheImpl) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	for {
		var keys []string
		var err error

		keys, cursor, err := r.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			r.client.Del(ctx, keys...)
		}

		if cursor == 0 {
			break
		}
	}
	return nil
}
