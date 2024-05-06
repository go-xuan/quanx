package cachex

import (
	"context"
	"time"
)

type CacheClient struct {
	key    func(string) string
	get    func(context.Context, string) any                 // 获取缓存值
	set    func(context.Context, string, any, time.Duration) // 设置缓存值
	delete func(context.Context, ...string) int64            // 删除缓存值
	exists func(context.Context, ...string) bool             // 是否存在
}

func GetValue[T any](ctx context.Context, cache *CacheClient, k string, get func() T) T {
	k = cache.key(k)
	var v T
	if cache.exists(ctx, k) {
		v = cache.get(ctx, k).(T)
	} else {
		v = get()
		cache.set(ctx, k, v, time.Duration(0))
	}
	return v
}

func (c *CacheClient) Get(ctx context.Context, k string) any {
	return c.get(ctx, c.key(k))
}

func (c *CacheClient) Set(ctx context.Context, k string, v any, d time.Duration) {
	c.set(ctx, c.key(k), v, d)
}

func (c *CacheClient) Del(ctx context.Context, k string) int64 {
	return c.delete(ctx, c.key(k))
}
