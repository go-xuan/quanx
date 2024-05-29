package cachex

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/types/slicex"
)

type Client interface {
	SET(ctx context.Context, key string, value any, expiration time.Duration)
	GET(ctx context.Context, key string) string
	DELETE(ctx context.Context, keys ...string) int64
	EXIST(ctx context.Context, keys ...string) bool
}

// 本地缓存客户端
type LocalClient struct {
	cache     *Cache
	client    *cache.Cache
	unmarshal func([]byte, any) error
}

func (c *LocalClient) SET(ctx context.Context, key string, value any, d time.Duration) {
	c.client.Set(c.cache.GetKey(key), value, d)
}

func (c *LocalClient) GET(ctx context.Context, key string) string {
	if v, ok := c.client.Get(c.cache.GetKey(key)); ok {
		return v.(string)
	}
	return ""
}

func (c *LocalClient) DELETE(ctx context.Context, keys ...string) int64 {
	if len(keys) > 0 {
		for _, key := range keys {
			c.client.Delete(c.cache.GetKey(key))
		}
		return int64(len(keys))
	}
	return 0
}

func (c *LocalClient) EXIST(ctx context.Context, keys ...string) bool {
	if len(keys) > 0 {
		for _, k := range keys {
			if _, ok := c.client.Get(c.cache.GetKey(k)); !ok {
				return false
			}
		}
		return true
	}
	return false
}

// redis缓存客户端
type RedisClient struct {
	cache     *Cache
	client    redis.UniversalClient
	unmarshal func([]byte, any) error
}

func (c *RedisClient) SET(ctx context.Context, key string, value any, expiration time.Duration) {
	c.client.Set(ctx, c.cache.GetKey(key), value, expiration)
}

func (c *RedisClient) GET(ctx context.Context, key string) string {
	if value, err := c.client.Get(ctx, c.cache.GetKey(key)).Result(); err == nil {
		return value
	}
	return ""
}

func (c *RedisClient) DELETE(ctx context.Context, keys ...string) int64 {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var n int64
			if n, err = c.client.Del(ctx, c.cache.GetKeys(keys[x:y])...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return total
		}
	} else {
		if total, err = c.client.Del(ctx, c.cache.GetKeys(keys)...).Result(); err != nil {
			return total
		}
	}
	return 0
}

func (c *RedisClient) EXIST(ctx context.Context, keys ...string) bool {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var n int64
			if n, err = c.client.Exists(ctx, c.cache.GetKeys(keys[x:y])...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return false
		}
	} else {
		if total, err = c.client.Exists(ctx, c.cache.GetKeys(keys)...).Result(); err != nil {
			return false
		}
	}
	return total > 0
}
