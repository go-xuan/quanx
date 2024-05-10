package cachex

import (
	"fmt"
	"github.com/go-xuan/quanx/db/redisx"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

const (
	CacheTypeRedis = "redis"
	CacheTypeLocal = "local"
)

// 缓存配置
type MultiCache []*Cache

type Cache struct {
	Type    string `json:"type" yaml:"type" default:"redis"`         // 缓存类型（local/redis）
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
	if handler == nil {
		handler = &Handler{
			Multi:     false,
			Client:    client,
			ClientMap: make(map[string]*CacheClient),
		}
	} else {
		handler.Multi = true
	}
	handler.ClientMap[c.Source] = client
	log.Info("Cache Init Successful: ", c.ToString())
	return
}

func Default() *Cache {
	return &Cache{
		Source:  constx.DefaultSourceName,
		Prefix:  "cache",
		Marshal: marshalx.Json,
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
	if handler == nil {
		handler = &Handler{
			Multi:     true,
			ClientMap: make(map[string]*CacheClient),
		}
	} else {
		handler.Multi = true
	}
	multi := anyx.IfZero(m, MultiCache{Default()})
	for i, c := range multi {
		var client = c.CacheClient()
		handler.ClientMap[c.Source] = client
		if i == 0 || c.Source == constx.DefaultSourceName {
			handler.Client = client
		}
		log.Info("Cache Init Successful: ", c.ToString())
	}
	return nil
}

// 配置信息格式化
func (c *Cache) ToString() string {
	return fmt.Sprintf("type=%s source=%s prefix=%s marshal=%s", c.Type, c.Source, c.Prefix, c.Marshal)
}

func (c *Cache) CacheClient() *CacheClient {
	if c.Type == CacheTypeRedis {
		return &CacheClient{
			cache:     c,
			client:    &RedisClient{redisx.Client(c.Source)},
			unmarshal: marshalx.NewCase(c.Marshal).Unmarshal,
		}
	} else {
		return &CacheClient{
			cache:     c,
			client:    &LocalClient{localCache},
			unmarshal: marshalx.NewCase(c.Marshal).Unmarshal,
		}
	}
}

func (c *Cache) GetKey(key string) string {
	if c.Prefix != "" {
		return stringx.AddSuffix(c.Prefix, ":") + key
	} else {
		return key
	}
}

func (c *Cache) GetKeys(keys []string) []string {
	if len(keys) > 0 && c.Prefix != "" {
		prefix := stringx.AddSuffix(c.Prefix, ":")
		var newKeys []string
		for _, key := range keys {
			newKeys = append(newKeys, prefix+key)
		}
		return newKeys
	}
	return keys
}
