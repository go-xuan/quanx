package cachex

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/types/slicex"
)

type RedisClient struct {
	client redis.UniversalClient
}

func (c *RedisClient) SET(ctx context.Context, key string, value any, expiration time.Duration) {
	c.client.Set(ctx, key, value, expiration)
}

func (c *RedisClient) GET(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisClient) DELETE(ctx context.Context, keys ...string) int64 {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) error {
			var n int64
			if n, err = c.client.Del(ctx, keys[x:y]...).Result(); err != nil {
				return err
			}
			total += n
			return nil
		}); err != nil {
			return total
		}
	} else {
		if total, err = c.client.Del(ctx, keys...).Result(); err != nil {
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
			if n, err = c.client.Exists(ctx, keys[x:y]...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return false
		}
	} else {
		if total, err = c.client.Exists(ctx, keys...).Result(); err != nil {
			return false
		}
	}
	return total > 0
}
