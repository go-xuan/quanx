package cachex

import (
	"context"
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/patrickmn/go-cache"
)

// LocalClientBuilder 本地缓存客户端构建器
func LocalClientBuilder(config *Config) (Client, error) {
	return NewLocalClient(config)
}

// NewLocalClient 创建本地缓存客户端
func NewLocalClient(config *Config) (*LocalClient, error) {
	return &LocalClient{
		config:  config,
		marshal: marshalx.Apply(config.Marshal),
		cache:   cache.New(time.Duration(-1), time.Duration(-1)),
	}, nil
}

// LocalClient 本地缓存客户端
type LocalClient struct {
	config  *Config
	marshal marshalx.Marshal
	cache   *cache.Cache
}

func (c *LocalClient) GetCache() *cache.Cache {
	return c.cache
}

func (c *LocalClient) GetConfig() *Config {
	return c.config
}

func (c *LocalClient) GetInstance() interface{} {
	return c.cache
}

func (c *LocalClient) GetKey(key string) string {
	return c.GetConfig().GetKey(key)
}

func (c *LocalClient) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	bytes, err := c.marshal.Marshal(value)
	if err != nil {
		return errorx.Wrap(err, "marshal value error")
	}
	c.GetCache().Set(c.GetKey(key), string(bytes), expiration)
	return nil
}

func (c *LocalClient) Get(ctx context.Context, key string, value any) bool {
	if result := c.GetString(ctx, key); result != "" {
		if err := c.marshal.Unmarshal([]byte(result), value); err == nil {
			return true
		}
	}
	return false
}

func (c *LocalClient) GetString(_ context.Context, key string) string {
	if result, ok := c.GetCache().Get(c.GetKey(key)); ok {
		return result.(string)
	}
	return ""
}

func (c *LocalClient) Delete(_ context.Context, key string) bool {
	c.GetCache().Delete(c.GetKey(key))
	return true
}

func (c *LocalClient) Exist(_ context.Context, key string) bool {
	_, ok := c.GetCache().Get(c.GetKey(key))
	return ok
}

func (c *LocalClient) Expire(_ context.Context, key string, expiration time.Duration) error {
	key = c.GetKey(key)
	result, ok := c.GetCache().Get(key)
	if !ok {
		return errorx.New("key not found")
	}
	c.GetCache().Set(key, result, expiration)
	return nil
}

func (c *LocalClient) Copy(source string, database int) (Client, error) {
	return &LocalClient{
		config:  c.config,
		marshal: c.marshal,
		cache:   c.cache,
	}, nil
}
