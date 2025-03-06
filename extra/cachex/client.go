package cachex

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/taskx"
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
	marshal marshalx.Method
}

func (c *LocalClient) Config() *Config {
	return c.config
}

func (c *LocalClient) Set(ctx context.Context, key string, value any, d time.Duration) error {
	if bytes, err := c.marshal.Marshal(value); err != nil {
		return errorx.Wrap(err, "marshal value error")
	} else {
		c.client.Set(c.config.GetKey(key), string(bytes), d)
		return nil
	}
}

func (c *LocalClient) Get(ctx context.Context, key string, value any) bool {
	if result := c.GetString(ctx, key); result != "" {
		if err := c.marshal.Unmarshal([]byte(result), value); err == nil {
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
	if result, ok := c.client.Get(key); !ok {
		return errorx.New("key not found")
	} else {
		c.client.Set(key, result, d)
		return nil
	}
}

// RedisClient redis缓存客户端
type RedisClient struct {
	config  *Config
	client  redis.UniversalClient
	marshal marshalx.Method
}

func (c *RedisClient) Config() *Config {
	return c.config
}

func (c *RedisClient) Set(ctx context.Context, key string, value any, d time.Duration) error {
	if bytes, err := c.marshal.Marshal(value); err != nil {
		return errorx.Wrap(err, "marshal value error")
	} else if err = c.client.Set(ctx, c.config.GetKey(key), bytes, d).Err(); err != nil {
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
	if result, err := c.client.Get(ctx, c.config.GetKey(key)).Result(); err == nil {
		return result
	}
	return ""
}

func (c *RedisClient) Delete(ctx context.Context, keys ...string) int64 {
	var sum int64
	_ = taskx.ExecWithBatches(len(keys), 100, func(x int, y int) error {
		if s, err := c.client.Del(ctx, c.config.GetKeys(keys[x:y])...).Result(); err != nil {
			return errorx.Wrap(err, "delete redis keys error")
		} else {
			sum += s
			return nil
		}
	})
	return sum
}

func (c *RedisClient) Exist(ctx context.Context, keys ...string) bool {
	var sum int64
	_ = taskx.ExecWithBatches(len(keys), 100, func(x int, y int) error {
		if n, err := c.client.Exists(ctx, c.config.GetKeys(keys[x:y])...).Result(); err != nil {
			return errorx.Wrap(err, "redis exists error")
		} else {
			sum += n
			return nil
		}
	})
	return sum > 0
}

func (c *RedisClient) Expire(ctx context.Context, key string, d time.Duration) error {
	if err := c.client.Expire(ctx, c.config.GetKey(key), d).Err(); err != nil {
		return errorx.Wrap(err, "redis expire error")
	}
	return nil
}
