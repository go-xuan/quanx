package cachex

import (
	"context"
	"errors"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/go-xuan/quanx/core/confx"
)

var localCache *cache.Cache

type LocalCache struct {
	Expiration int    `json:"expiration" yaml:"expiration" default:"-1"` // 缓存有效期
	Interval   int    `json:"interval" yaml:"interval" default:"-1"`     // 缓存清理间隔
	File       string `json:"file" yaml:"file"`                          // 缓存加载文件
}

// 配置信息格式化
func (LocalCache) Theme() string {
	return "Redis"
}

// 配置文件读取
func (LocalCache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "local_cache.yaml",
		NacosDataId: "local_cache.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (conf LocalCache) Run() error {
	localCache = cache.New(time.Duration(conf.Expiration), time.Duration(conf.Interval))
	if conf.File != "" {
		return localCache.LoadFile(conf.File)
	}
	return nil
}

type LocalClient struct {
	cache *cache.Cache
}

func (c *LocalClient) SET(ctx context.Context, k string, v any, d time.Duration) {
	c.cache.Set(k, v, d)
}

func (c *LocalClient) GET(ctx context.Context, k string) (string, error) {
	if v, ok := c.cache.Get(k); ok {
		return v.(string), nil
	}
	return "", errors.New("key not found")
}

func (c *LocalClient) DELETE(ctx context.Context, keys ...string) int64 {
	if len(keys) > 0 {
		for _, k := range keys {
			c.cache.Delete(k)
		}
		return int64(len(keys))
	}
	return 0
}

func (c *LocalClient) EXIST(ctx context.Context, keys ...string) bool {
	if len(keys) > 0 {
		for _, k := range keys {
			if _, ok := c.cache.Get(k); !ok {
				return false
			}
		}
		return true
	}
	return false
}
