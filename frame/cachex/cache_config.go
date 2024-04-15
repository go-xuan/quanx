package cachex

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/frame/confx"
	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/go-xuan/quanx/utils/anyx"
	"github.com/go-xuan/quanx/utils/slicex"
)

// 缓存配置
type MultiCache []*Cache

type Cache struct {
	Source      string         `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix      string         `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal     string         `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
	config      *redisx.Redis  // redis配置
	cmdable     redis.Cmdable  // redis连接
	marshalCase *marshalx.Case // 序列化方案
}

func (c *Cache) Theme() string {
	return "Cache"
}

func (c *Cache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (c *Cache) Run() (err error) {
	if err = anyx.SetDefaultValue(c); err != nil {
		return
	}
	// 初始化缓存属性
	c.InitAttribute()
	client := c.CacheClient()
	handler = &Handler{
		Multi:     false,
		Client:    client,
		Config:    c,
		ClientMap: make(map[string]*CacheClient[any]),
		ConfigMap: make(map[string]*Cache),
	}
	handler.ClientMap[c.Source] = client
	handler.ConfigMap[c.Source] = c
	log.Info(c.ToString("Cache init successful!"))
	return
}

func Default() *Cache {
	return &Cache{
		Source:  constx.Default,
		Prefix:  "cache",
		Marshal: marshalx.Msgpack,
	}
}

// 配置信息格式化
func (MultiCache) Theme() string {
	return "Cache"
}

// 配置文件读取
func (MultiCache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

// 配置器运行
func (m MultiCache) Run() error {
	handler = &Handler{
		Multi:     true,
		ClientMap: make(map[string]*CacheClient[any]),
		ConfigMap: make(map[string]*Cache),
	}
	multi := anyx.IfZero(m, MultiCache{Default()})
	for i, cache := range multi {
		// 初始化缓存属性
		cache.InitAttribute()
		client := cache.CacheClient()
		handler.ClientMap[cache.Source] = client
		handler.ConfigMap[cache.Source] = cache
		if i == 0 || cache.Source == constx.Default {
			handler.Client = client
			handler.Config = cache
		}
		log.Info(cache.ToString("init cache successful!"))
	}
	return nil
}

// 配置信息格式化
func (c *Cache) ToString(title string) string {
	return fmt.Sprintf("%s => source=%s prefix=%s marshal=%s",
		title, c.Source, c.Prefix, c.Marshal)
}

// 初始化属性
func (c *Cache) InitAttribute() {
	if c.config == nil {
		c.config = redisx.This().GetConfig(c.Source)
	}
	if c.cmdable == nil {
		c.cmdable = redisx.This().GetCmdable(c.Source)
	}
	if c.marshalCase == nil {
		c.marshalCase = marshalx.DefaultCase()
	}
}

func (c *Cache) CacheClient() *CacheClient[any] {
	return &CacheClient[any]{
		Set: func(ctx context.Context, id string, v any, duration time.Duration) {
			_ = c.Set(ctx, id, v, duration)
		},
		Get: func(ctx context.Context, id string) (v any) {
			if exist, err := c.Exists(ctx, id); exist > 0 && err == nil {
				if err = c.Get(ctx, id, v); err != nil {
					return
				}
			}
			return
		},
		Del: func(ctx context.Context, s string) {
			_, _ = c.Delete(ctx, s)
		},
		Exist: func(ctx context.Context, s string) bool {
			if exist, err := c.Exists(ctx, s); exist > 0 && err == nil {
				return true
			}
			return false
		},
	}
}

func (c *Cache) CASE() *marshalx.Case {
	if c.marshalCase == nil {
		c.marshalCase = marshalx.DefaultCase()
	}
	return c.marshalCase
}

func (c *Cache) RedisDB() redis.Cmdable {
	if c.cmdable == nil {
		c.cmdable = redisx.This().GetCmdable(c.Source)
	}
	return c.cmdable
}

func (c *Cache) Get(ctx context.Context, key string, v any) error {
	if value, err := c.RedisDB().Get(ctx, c.Prefix+key).Bytes(); err != nil {
		return err
	} else {
		return c.CASE().Unmarshal(value, v)
	}
}

func (c *Cache) Set(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := c.CASE().Marshal(v); err != nil {
		return err
	} else {
		return c.RedisDB().Set(ctx, c.Prefix+key, value, expiration).Err()
	}
}

func (c *Cache) SetNX(ctx context.Context, key string, v any, expiration time.Duration) error {
	if value, err := c.CASE().Marshal(v); err != nil {
		return err
	} else {
		return c.RedisDB().SetNX(ctx, c.Prefix+key, value, expiration).Err()
	}
}

func (c *Cache) Exists(ctx context.Context, keys ...string) (total int64, err error) {
	if l := len(keys); l > 0 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var batches = c.AddPrefix(keys[x:y])
			var n int64
			if n, err = c.RedisDB().Exists(ctx, batches...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return
		}
	}
	return
}

func (c *Cache) Delete(ctx context.Context, keys ...string) (total int64, err error) {
	if l := len(keys); l > 0 {
		if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
			var batches = c.AddPrefix(keys[x:y])
			var n int64
			if n, err = c.RedisDB().Del(ctx, batches...).Result(); err != nil {
				return
			}
			total += n
			return
		}); err != nil {
			return
		}
	}
	return
}

func (c *Cache) AddPrefix(keys []string) []string {
	var newKeys []string
	if len(keys) > 0 {
		for _, key := range keys {
			newKeys = append(newKeys, c.Prefix+key)
		}
	}
	return newKeys
}
