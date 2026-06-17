package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jereminathanael/healthqueue/internal/config"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func New(cfg *config.Config) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr(),
		Password: cfg.RedisPassword,
		DB: cfg.RedisDB,
		PoolSize: 10,
		MinIdleConns: 3,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Cache{client: client}, nil
}

// Set simpan value apapun ke JSON dengan TTL
func (c *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get ambil cached value. Return (false, nil) kalau cache miss — bukan error!
func (c *Cache) Get(ctx context.Context, key string, dest any) (bool, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil //cache miss
	}
	if err != nil {
		return false, fmt.Errorf("cache get error: %w", err)
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, fmt.Errorf("unmarshall error: %w", err)
	}
	return true, nil
}

// Delete hapus satu atau lebih key
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	return c.client.Del(ctx, keys...).Err()
}

// DeleteByPattern hapus semua key yang match pattern, misal "doctors:*"
func (c *Cache) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...). Err()
}

func (c *Cache) Close() error {
	return c.client.Close()
}

