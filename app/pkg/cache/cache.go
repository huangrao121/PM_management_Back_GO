package cache

import (
	"context"
	"encoding/json"
	"time"

	"pm_go_version/app/pkg/redis"

	redisClient "github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redisClient.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: redis.GetClient(),
	}
}

func (c *RedisCache) Get(ctx context.Context, key string, value interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), value)
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
