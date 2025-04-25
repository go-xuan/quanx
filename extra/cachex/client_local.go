package cachex

import (
	"context"
	"time"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/utils/marshalx"
	"github.com/patrickmn/go-cache"
)

// LocalClient 本地缓存客户端
type LocalClient struct {
	config  *Config
	cache   *cache.Cache
	marshal marshalx.Method
}

func (c *LocalClient) Instance() *cache.Cache {
	return c.cache
}

func (c *LocalClient) Config() *Config {
	return c.config
}

func (c *LocalClient) Set(ctx context.Context, key string, value any, d time.Duration) error {
	if bytes, err := c.marshal.Marshal(value); err != nil {
		return errorx.Wrap(err, "marshal value error")
	} else {
		c.Instance().Set(c.config.GetKey(key), string(bytes), d)
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
	if result, ok := c.Instance().Get(c.config.GetKey(key)); ok {
		return result.(string)
	}
	return ""
}

func (c *LocalClient) Delete(ctx context.Context, keys ...string) int64 {
	if len(keys) > 0 {
		for _, key := range keys {
			c.Instance().Delete(c.config.GetKey(key))
		}
		return int64(len(keys))
	}
	return 0
}

func (c *LocalClient) Exist(ctx context.Context, keys ...string) bool {
	if len(keys) > 0 {
		for _, k := range keys {
			if _, ok := c.Instance().Get(c.config.GetKey(k)); !ok {
				return false
			}
		}
		return true
	}
	return false
}

func (c *LocalClient) Expire(ctx context.Context, key string, d time.Duration) error {
	key = c.config.GetKey(key)
	if result, ok := c.Instance().Get(key); !ok {
		return errorx.New("key not found")
	} else {
		c.Instance().Set(key, result, d)
		return nil
	}
}
