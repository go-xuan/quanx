package cachex

import (
	"context"
	"github.com/go-xuan/quanx/types/slicex"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheClient struct {
	cache     *Cache
	client    redis.UniversalClient
	unmarshal func([]byte, any) error
}

func GetValue[T any](ctx context.Context, client *CacheClient, k string, get func() T) T {
	k = client.cache.Prefix + k
	var v T
	if client.Exists(ctx, k) {
		client.Get(ctx, k, v)
	} else {
		v = get()
		client.Set(ctx, k, v, time.Duration(0))
	}
	return v
}

func (c *CacheClient) Get(ctx context.Context, k string, v any) {
	if value, err := c.client.Get(ctx, c.cache.Prefix+k).Bytes(); err == nil {
		_ = c.unmarshal(value, v)
	}
	return
}

func (c *CacheClient) GetString(ctx context.Context, k string) string {
	return c.client.Get(ctx, c.cache.Prefix+k).Val()
}

func (c *CacheClient) GetInt(ctx context.Context, k string) (int, error) {
	return c.client.Get(ctx, c.cache.Prefix+k).Int()
}

func (c *CacheClient) GetInt64(ctx context.Context, k string) (int64, error) {
	return c.client.Get(ctx, c.cache.Prefix+k).Int64()
}

func (c *CacheClient) GetTime(ctx context.Context, k string) (time.Time, error) {
	return c.client.Get(ctx, c.cache.Prefix+k).Time()
}

func (c *CacheClient) Set(ctx context.Context, k string, v any, d time.Duration) {
	c.client.Set(ctx, c.cache.Prefix+k, v, d)
}

func (c *CacheClient) Del(ctx context.Context, k ...string) int64 {
	var total int64
	var err error
	if l := len(k); l > 100 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) error {
			var keys []string
			for _, key := range k[x:y] {
				keys = append(keys, c.cache.Prefix+key)
			}
			var n int64
			if n, err = c.client.Del(ctx, keys...).Result(); err != nil {
				return err
			}
			total += n
			return nil
		}); err != nil {
			return total
		}
	} else {
		var keys []string
		for _, key := range k {
			keys = append(keys, c.cache.Prefix+key)
		}
		if total, err = c.client.Del(ctx, keys...).Result(); err != nil {
			return total
		}
	}
	return 0
}

func (c *CacheClient) Exists(ctx context.Context, k ...string) bool {
	var total int64
	var err error
	if l := len(k); l > 100 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var keys []string
			for _, key := range k[x:y] {
				keys = append(keys, c.cache.Prefix+key)
			}
			var n int64
			if n, err = c.client.Exists(ctx, keys...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return false
		}
	} else {
		var keys []string
		for _, key := range k {
			keys = append(keys, c.cache.Prefix+key)
		}
		if total, err = c.client.Exists(ctx, keys...).Result(); err != nil {
			return false
		}
	}
	return total > 0
}
