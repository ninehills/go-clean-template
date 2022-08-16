package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

const DefaultCacheExpires = 5 * time.Minute

var (
	ErrMiss      = errors.New("cache miss")
	ErrMarshal   = errors.New("cache marshal error")
	ErrUnmarshal = errors.New("cache unmarshal error")
	ErrStorage   = errors.New("cache storage visit error")
)

type Cache struct {
	redis   *redis.Client
	expires time.Duration
}

func NewCache(r *redis.Client, expires time.Duration) *Cache {
	return &Cache{
		redis:   r,
		expires: expires,
	}
}

// 从 Cache 中获取对象.
func (c *Cache) Get(ctx context.Context, key string, target interface{}) error {
	result, err := c.redis.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return ErrMiss
	} else if err != nil {
		return ErrStorage
	}

	err = json.Unmarshal(result, target)
	if err != nil {
		return ErrUnmarshal
	}

	return nil
}

// 写入 Cache.
func (c *Cache) Set(ctx context.Context, key string, value interface{}) error {
	result, err := json.Marshal(value)
	if err != nil {
		return ErrMarshal
	}

	err = c.redis.Set(ctx, key, result, c.expires).Err()
	if err != nil {
		return ErrStorage
	}

	return nil
}

// 删除 Cache.
func (c *Cache) Del(ctx context.Context, key string) error {
	err := c.redis.Del(ctx, key).Err()
	if errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		return ErrStorage
	}

	return nil
}
