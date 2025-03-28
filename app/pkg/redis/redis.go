package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

// InitRedis 初始化 Redis 客户端
func InitRedis(addr, password string, db int) error {
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis connection failed: %v", err)
	}

	return nil
}

// GetClient 获取 Redis 客户端实例
func GetClient() *redis.Client {
	return client
}
