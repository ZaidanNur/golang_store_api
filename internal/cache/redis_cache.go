package cache

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr string) *RedisCache {
	if addr == "" {
		addr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("[CACHE] WARNING: Redis not reachable at %s: %v — caching disabled", addr, err)
		return &RedisCache{client: nil}
	}

	log.Printf("[CACHE] Connected to Redis at %s", addr)
	return &RedisCache{client: client}
}

func (c *RedisCache) IsAvailable() bool {
	return c.client != nil
}

func (c *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	if c.client == nil {
		return nil, nil
	}
	val, err := c.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (c *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	if c.client == nil {
		return nil
	}
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if c.client == nil {
		return nil
	}
	return c.client.Del(ctx, key).Err()
}

func (c *RedisCache) DeleteByPrefix(ctx context.Context, prefix string) error {
	if c.client == nil {
		return nil
	}
	iter := c.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
	return iter.Err()
}
