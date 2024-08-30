package cachex

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/confx"
	"github.com/go-xuan/quanx/app/constx"
	"github.com/go-xuan/quanx/server/redisx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

const (
	CacheTypeRedis = "redis"
	CacheTypeLocal = "local"
)

// MultiCache 缓存配置
type MultiCache []*Cache

type Cache struct {
	Type    string `json:"type" yaml:"type" default:"redis"`         // 缓存类型（local/redis）
	Source  string `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
}

// Title 配置器标题
func (c *Cache) Title() string {
	return "Cache"
}

// Reader 配置文件读取
func (c *Cache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (c *Cache) Run() (err error) {
	if err = anyx.SetDefaultValue(c); err != nil {
		return
	}
	var client = c.InitClient()
	if handler == nil {
		handler = &Handler{
			Multi:     false,
			Client:    client,
			ClientMap: make(map[string]Client),
		}
	} else {
		handler.Multi = true
	}
	handler.ClientMap[c.Source] = client
	log.Info("Cache Init Successful: ", c.Info())
	return
}

func Default() *Cache {
	return &Cache{
		Source:  constx.DefaultKey,
		Prefix:  "cache",
		Marshal: "json",
	}
}

// Title 配置信息格式化
func (MultiCache) Title() string {
	return "Cache"
}

// Reader 配置文件读取
func (MultiCache) Reader() *confx.Reader {
	return &confx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

// Run 配置器运行
func (m MultiCache) Run() error {
	if handler == nil {
		handler = &Handler{
			Multi:     true,
			ClientMap: make(map[string]Client),
		}
	} else {
		handler.Multi = true
	}
	multi := anyx.IfZero(m, MultiCache{Default()})
	for i, c := range multi {
		var client = c.InitClient()
		handler.ClientMap[c.Source] = client
		if i == 0 || c.Source == constx.DefaultKey {
			handler.Client = client
		}
		log.Info("Cache Init Successful: ", c.Info())
	}
	return nil
}

// Info 配置信息格式化
func (c *Cache) Info() string {
	return fmt.Sprintf("type=%s source=%s prefix=%s marshal=%s", c.Type, c.Source, c.Prefix, c.Marshal)
}

// InitClient 根据缓存配置初始化缓存客户端
func (c *Cache) InitClient() Client {
	switch c.Type {
	case CacheTypeRedis:
		return &RedisClient{
			cache:   c,
			client:  redisx.Client(c.Source),
			marshal: marshalx.NewCase(c.Marshal),
		}
	case CacheTypeLocal:
		return &LocalClient{
			cache:   c,
			client:  NewLocalCache(),
			convert: marshalx.NewCase(c.Marshal),
		}
	default:
		log.Error("cache client not support type: ", c.Type)
		return nil
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
