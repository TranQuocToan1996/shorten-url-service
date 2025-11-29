package redis_cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"shorten/pkg/cache"
)

type redisCache struct {
	client *redis.Client
}

var _ cache.Cache = (*redisCache)(nil)

func NewRedisCache(client *redis.Client) cache.Cache {
	return &redisCache{
		client: client,
	}
}

func (r *redisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, cache.ErrNotFound
	}
	return val, err
}

func (r *redisCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *redisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

