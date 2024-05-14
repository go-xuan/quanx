package cachex

import (
	"context"
	"time"
)

type Client interface {
	SET(ctx context.Context, k string, v any, expiration time.Duration)
	GET(ctx context.Context, k string) (string, error)
	DELETE(ctx context.Context, k ...string) int64
	EXIST(ctx context.Context, k ...string) bool
}

type CacheClient struct {
	cache     *Cache
	client    Client
	unmarshal func([]byte, any) error
}

func GetValue[T any](ctx context.Context, client *CacheClient, k string, get func() T) T {
	k = client.cache.GetKey(k)
	var v T
	if client.Exists(ctx, k) {
		client.Get(ctx, k, v)
	} else {
		v = get()
		client.Set(ctx, k, v, time.Duration(0))
	}
	return v
}

func (c *CacheClient) Get(ctx context.Context, key string, value any) {
	if v, err := c.client.GET(ctx, c.cache.GetKey(key)); err == nil {
		_ = c.unmarshal([]byte(v), value)
	}
	return
}

func (c *CacheClient) GetString(ctx context.Context, key string) string {
	res, _ := c.client.GET(ctx, c.cache.GetKey(key))
	return res
}

func (c *CacheClient) Set(ctx context.Context, key string, value any, expiration time.Duration) {
	c.client.SET(ctx, c.cache.GetKey(key), value, expiration)
}

func (c *CacheClient) Delete(ctx context.Context, keys ...string) int64 {
	keys = c.cache.GetKeys(keys)
	return c.client.DELETE(ctx, keys...)
}

func (c *CacheClient) Exists(ctx context.Context, keys ...string) bool {
	keys = c.cache.GetKeys(keys)
	return c.client.EXIST(ctx, keys...)
}
