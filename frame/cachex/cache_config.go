package cachex

import (
	"context"
	"fmt"
	"github.com/go-xuan/quanx/utils/slicex"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/frame/confx"
	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/go-xuan/quanx/utils/anyx"
)

// 缓存配置
type MultiCache []*Cache

type Cache struct {
	Source  string `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
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
	var client = c.CacheClient()
	handler = &Handler{
		Multi:     false,
		Client:    client,
		ClientMap: make(map[string]*CacheClient),
	}
	handler.ClientMap[c.Source] = client
	log.Info(c.ToString("Cache Init Successful!"))
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
		ClientMap: make(map[string]*CacheClient),
	}
	multi := anyx.IfZero(m, MultiCache{Default()})
	for i, cache := range multi {
		var client = cache.CacheClient()
		handler.ClientMap[cache.Source] = client
		if i == 0 || cache.Source == constx.Default {
			handler.Client = client
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

func (c *Cache) CacheClient() *CacheClient {
	var client = redisx.Client(c.Source)
	unmarshal := marshalx.NewCase(c.Marshal).Unmarshal
	return &CacheClient{
		key: func(k string) string {
			return c.Prefix + k
		},
		set: func(ctx context.Context, k string, v any, duration time.Duration) {
			client.Set(ctx, k, v, duration)
		},
		get: func(ctx context.Context, k string) (v any) {
			if value, err := client.Get(ctx, k).Bytes(); err == nil {
				_ = unmarshal(value, v)
			}
			return
		},
		delete: func(ctx context.Context, k ...string) int64 {
			var total int64
			var err error
			if l := len(k); l > 100 {
				if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
					var n int64
					if n, err = client.Del(ctx, k[x:y]...).Result(); err != nil {
						return
					}
					total += n
					return
				}); err != nil {
					return total
				}
			} else {
				if total, err = client.Del(ctx, k...).Result(); err != nil {
					return total
				}
			}
			return 0
		},
		exists: func(ctx context.Context, k ...string) bool {
			var total int64
			var err error
			if l := len(k); l > 100 {
				if err = slicex.ExecInBatches(l, 100, func(x int, y int) (err error) {
					var n int64
					if n, err = client.Exists(ctx, k[x:y]...).Result(); err != nil {
						return
					}
					total += n
					return
				}); err != nil {
					return false
				}
			} else {
				if total, err = client.Exists(ctx, k...).Result(); err != nil {
					return false
				}
			}
			return total > 0
		},
	}
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
