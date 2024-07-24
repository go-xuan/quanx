package cachex

import (
	"context"
	"github.com/go-xuan/quanx/os/execx"
	"github.com/go-xuan/quanx/utils/marshalx"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
)

type Client interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, value any)
	GetString(ctx context.Context, key string) string
	Delete(ctx context.Context, keys ...string) int64
	Exist(ctx context.Context, keys ...string) bool
}

// LocalClient 本地缓存客户端
type LocalClient struct {
	cache   *Cache
	client  *cache.Cache
	convert *marshalx.Case
}

func (c *LocalClient) Set(ctx context.Context, key string, value any, d time.Duration) (err error) {
	var bytes []byte
	if bytes, err = c.convert.Marshal(value); err != nil {
		return
	}
	c.client.Set(c.cache.GetKey(key), string(bytes), d)
	return
}

func (c *LocalClient) Get(ctx context.Context, key string, value any) {
	result := c.GetString(ctx, key)
	_ = c.convert.Unmarshal([]byte(result), value)
}

func (c *LocalClient) GetString(ctx context.Context, key string) string {
	if result, ok := c.client.Get(c.cache.GetKey(key)); ok {
		return result.(string)
	}
	return ""
}

func (c *LocalClient) Delete(ctx context.Context, keys ...string) int64 {
	if len(keys) > 0 {
		for _, key := range keys {
			c.client.Delete(c.cache.GetKey(key))
		}
		return int64(len(keys))
	}
	return 0
}

func (c *LocalClient) Exist(ctx context.Context, keys ...string) bool {
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

// RedisClient redis缓存客户端
type RedisClient struct {
	cache   *Cache
	client  redis.UniversalClient
	convert *marshalx.Case
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) (err error) {
	var bytes []byte
	if bytes, err = c.convert.Marshal(value); err != nil {
		return
	}
	if err = c.client.Set(ctx, c.cache.GetKey(key), bytes, expiration).Err(); err != nil {
		return
	}
	return
}

func (c *RedisClient) Get(ctx context.Context, key string, value any) {
	result := c.GetString(ctx, key)
	_ = c.convert.Unmarshal([]byte(result), value)
}

func (c *RedisClient) GetString(ctx context.Context, key string) string {
	if result, err := c.client.Get(ctx, c.cache.GetKey(key)).Result(); err == nil {
		return result
	}
	return ""
}

func (c *RedisClient) Delete(ctx context.Context, keys ...string) int64 {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = execx.UseBatches(l, 100, func(x int, y int) (err error) {
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

func (c *RedisClient) Exist(ctx context.Context, keys ...string) bool {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = execx.UseBatches(l, 100, func(x int, y int) (err error) {
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
