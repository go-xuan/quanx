package cachex

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/execx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

type Client interface {
	Config() *Config                                                                // 获取配置
	Set(ctx context.Context, key string, value any, expiration time.Duration) error // 更新缓存
	Get(ctx context.Context, key string, value any) bool                            // 获取缓存（指针，任意类型）
	GetString(ctx context.Context, key string) string                               // 获取缓存（字符串类型）
	Expire(ctx context.Context, key string, expiration time.Duration) error         // 续期
	Delete(ctx context.Context, keys ...string) int64                               // 删除
	Exist(ctx context.Context, keys ...string) bool                                 // 是否存在
}

// LocalClient 本地缓存客户端
type LocalClient struct {
	config  *Config
	client  *cache.Cache
	convert *marshalx.Case
}

func (c *LocalClient) Config() *Config {
	return c.config
}

func (c *LocalClient) Set(ctx context.Context, key string, value any, d time.Duration) (err error) {
	var bytes []byte
	if bytes, err = c.convert.Marshal(value); err != nil {
		return
	}
	c.client.Set(c.config.GetKey(key), string(bytes), d)
	return
}

func (c *LocalClient) Get(ctx context.Context, key string, value any) bool {
	if result := c.GetString(ctx, key); result != "" {
		if err := c.convert.Unmarshal([]byte(result), value); err == nil {
			return true
		}
	}
	return false
}

func (c *LocalClient) GetString(ctx context.Context, key string) string {
	if result, ok := c.client.Get(c.config.GetKey(key)); ok {
		return result.(string)
	}
	return ""
}

func (c *LocalClient) Delete(ctx context.Context, keys ...string) int64 {
	if len(keys) > 0 {
		for _, key := range keys {
			c.client.Delete(c.config.GetKey(key))
		}
		return int64(len(keys))
	}
	return 0
}

func (c *LocalClient) Exist(ctx context.Context, keys ...string) bool {
	if len(keys) > 0 {
		for _, k := range keys {
			if _, ok := c.client.Get(c.config.GetKey(k)); !ok {
				return false
			}
		}
		return true
	}
	return false
}

func (c *LocalClient) Expire(ctx context.Context, key string, d time.Duration) error {
	key = c.config.GetKey(key)
	if result, ok := c.client.Get(key); ok {
		c.client.Set(key, result, d)
	}
	return errorx.New("")
}

// RedisClient redis缓存客户端
type RedisClient struct {
	config  *Config
	client  redis.UniversalClient
	marshal *marshalx.Case
}

func (c *RedisClient) Config() *Config {
	return c.config
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, d time.Duration) (err error) {
	var bytes []byte
	if bytes, err = c.marshal.Marshal(value); err != nil {
		return
	}
	if err = c.client.Set(ctx, c.config.GetKey(key), bytes, d).Err(); err != nil {
		return
	}
	return
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
	if result, err := c.client.Get(ctx, c.config.GetKey(key)).Result(); err == nil {
		return result
	}
	return ""
}

func (c *RedisClient) Delete(ctx context.Context, keys ...string) int64 {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = execx.InBatches(l, 100, func(x int, y int) (err error) {
			var n int64
			if n, err = c.client.Del(ctx, c.config.GetKeys(keys[x:y])...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return total
		}
	} else {
		if total, err = c.client.Del(ctx, c.config.GetKeys(keys)...).Result(); err != nil {
			return total
		}
	}
	return 0
}

func (c *RedisClient) Exist(ctx context.Context, keys ...string) bool {
	var total int64
	var err error
	if l := len(keys); l > 100 {
		if err = execx.InBatches(l, 100, func(x int, y int) (err error) {
			var n int64
			if n, err = c.client.Exists(ctx, c.config.GetKeys(keys[x:y])...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return false
		}
	} else {
		if total, err = c.client.Exists(ctx, c.config.GetKeys(keys)...).Result(); err != nil {
			return false
		}
	}
	return total > 0
}

func (c *RedisClient) Expire(ctx context.Context, key string, d time.Duration) error {
	if err := c.client.Expire(ctx, c.config.GetKey(key), d).Err(); err != nil {
		return err
	}
	return nil
}
