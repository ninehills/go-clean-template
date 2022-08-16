package cache

import "context"

type Cacher interface {
	// 从 Cache 中获取对象
	Get(ctx context.Context, key string, target interface{}) error
	// 写入 Cache
	Set(ctx context.Context, key string, value interface{}) error
	// 删除 Cache
	Del(ctx context.Context, key string) error
}
