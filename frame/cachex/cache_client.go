package cachex

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	config      *redisx.Redis                                     // redis配置
	cache       *Cache                                            // 缓存配置
	marshalCase *marshalx.Case                                    // 序列化方案
	client      redis.UniversalClient                             // redis连接客户端
	get         func(context.Context, string) any                 // 获取缓存值
	set         func(context.Context, string, any, time.Duration) // 设置缓存值
	delete      func(context.Context, ...string) int64            // 删除缓存值
	exists      func(context.Context, ...string) bool             // 是否存在
}

func (c *CacheClient) TotalKey(k string) string {
	return c.cache.Prefix + k
}

func (c *CacheClient) TotalKeys(k []string) []string {
	var newKeys []string
	if len(k) > 0 {
		for _, key := range k {
			newKeys = append(newKeys, c.cache.Prefix+key)
		}
	}
	return newKeys
}

func (c *CacheClient) GetValue(ctx context.Context, key string, get func() any) any {
	if c.exists(ctx, key) {
		return c.get(ctx, key)
	} else {
		v := get()
		c.set(ctx, key, v, time.Duration(0))
		return v
	}
}

func (c *CacheClient) Get(ctx context.Context, k string) any {
	return c.get(ctx, c.TotalKey(k))
}

func (c *CacheClient) Set(ctx context.Context, k string, v any, d time.Duration) {
	c.set(ctx, c.TotalKey(k), v, d)
}

func (c *CacheClient) Del(ctx context.Context, k string) int64 {
	return c.delete(ctx, c.TotalKey(k))
}
