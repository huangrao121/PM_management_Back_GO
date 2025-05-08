package redis_config

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var rc *RedisCache
var isRedisAvailable bool

func InitRedis() error {
	rc = NewRedisCache()
	ctx := context.Background()
	_, err := rc.client.Ping(ctx).Result()
	if err != nil {
		isRedisAvailable = false
		log.Error("Redis connection failed: ", err)
		return fmt.Errorf("redis连接失败: %v", err)
	}
	isRedisAvailable = true
	log.Info("Redis connection successful")
	return nil
}

func GetRedisCache() *RedisCache {
	if !isRedisAvailable {
		return &RedisCache{
			client: nil,
		}
	}
	return rc
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache() *RedisCache {
	return &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func CloseRedis() {
	rc.client.Close()
}

func (rc *RedisCache) acquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	lockKey := fmt.Sprintf("lock%s", key)
	return rc.client.SetNX(ctx, lockKey, "locked", expiration).Result()
}

func (rc *RedisCache) releaseLock(ctx context.Context, key string) error {
	lockKey := fmt.Sprintf("lock%s", key)
	return rc.client.Del(ctx, lockKey).Err()
}

func (rc *RedisCache) GetStructValue(
	ctx context.Context,
	key string,
	loadFunc func() (interface{}, error),
) (string, error) {
	if !isRedisAvailable || rc.client == nil {
		log.Warn("Redis is not available, falling back to direct function call")
		value, err := loadFunc()
		if err != nil {
			return "", err
		}
		jsonData, err := json.Marshal(value)
		log.Debug("jsonData: ", string(jsonData))
		if err != nil {
			return "", fmt.Errorf("failed to serialize: %v", err)
		}
		return string(jsonData), nil
	}

	log.Info("Trying to get value from Redis with key: ", key)
	data, err := rc.client.Get(ctx, key).Result()
	if err == nil {
		log.Info("Successfully retrieved data from Redis, data length: ", len(data))
		return data, nil
	}

	log.Info("Data not found in Redis, acquiring lock")
	locked, err := rc.acquireLock(ctx, key, time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to get the lock %v", err)
	}

	if !locked {
		log.Info("Lock acquisition failed, retrying after delay")
		time.Sleep(60 * time.Millisecond)
		return rc.GetStructValue(ctx, key, loadFunc)
	}

	defer func() {
		_ = rc.releaseLock(ctx, key)
	}()

	log.Info("Loading data from function")
	value, err := loadFunc()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to serialize: %v", err)
	}

	log.Info("Storing data in Redis with key: ", key, ", data length: ", len(jsonData))
	err = rc.client.Set(ctx, key, jsonData, 30*time.Second).Err()
	if err != nil {
		log.Error("Failed to store data in Redis: ", err)
		return "", err
	}

	// 验证数据是否真的存储成功
	storedData, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		log.Error("Failed to verify stored data: ", err)
	} else {
		log.Info("Verified stored data length: ", len(storedData))
	}

	return string(jsonData), nil
}

// var RedisClient *redis.Client

// func InitRedis() error {
// 	RedisClient = redis.NewClient(&redis.Options{
// 		Addr:     os.Getenv("REDIS_ADDR"),
// 		Password: os.Getenv("REDIS_PASSWORD"),
// 		DB:       0, // 使用默认 DB
// 	})

// 	// 测试连接
// 	ctx := context.Background()
// 	_, err := RedisClient.Ping(ctx).Result()
// 	if err != nil {
// 		return fmt.Errorf("redis连接失败: %v", err)
// 	}

// 	return nil
// }

// // GetRedisClient 获取 Redis 客户端实例
// func GetRedisClient() *redis.Client {
// 	return RedisClient
// }
