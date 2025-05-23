package cachex

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// RedisClient redis缓存客户端
type RedisClient struct {
	config  *Config
	client  redis.UniversalClient
	marshal marshalx.Method
}

func (c *RedisClient) Instance() redis.UniversalClient {
	return c.client
}

func (c *RedisClient) Config() *Config {
	return c.config
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, d time.Duration) error {
	if bytes, err := c.marshal.Marshal(value); err != nil {
		return errorx.Wrap(err, "marshal value error")
	} else if err = c.Instance().Set(ctx, c.config.GetKey(key), bytes, d).Err(); err != nil {
		return errorx.Wrap(err, "set value error")
	}
	return nil
}

func (c *RedisClient) Get(ctx context.Context, key string, value any) bool {
	if result := c.GetString(ctx, key); result != "" {
		if err := c.marshal.Unmarshal([]byte(result), value); err == nil {
			return true
		}
	}
	return false

}

func (c *RedisClient) GetString(ctx context.Context, key string) string {
	if result, err := c.Instance().Get(ctx, c.config.GetKey(key)).Result(); err == nil {
		return result
	}
	return ""
}

func (c *RedisClient) Delete(ctx context.Context, key string) bool {
	if n, err := c.Instance().Del(ctx, c.config.GetKey(key)).Result(); err == nil && n > 0 {
		return true
	}
	return false
}

func (c *RedisClient) Exist(ctx context.Context, key string) bool {
	if n, err := c.Instance().Exists(ctx, c.config.GetKey(key)).Result(); err == nil && n > 0 {
		return true
	}
	return false
}

func (c *RedisClient) Expire(ctx context.Context, key string, d time.Duration) error {
	if err := c.Instance().Expire(ctx, c.config.GetKey(key), d).Err(); err != nil {
		return errorx.Wrap(err, "redis expire error")
	}
	return nil
}
