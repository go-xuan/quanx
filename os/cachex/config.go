package cachex

import (
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/app/configx"
	"github.com/go-xuan/quanx/app/constx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/server/redisx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/fmtx"
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

func (c *Cache) ID() string {
	return "cache"
}

func (c *Cache) Format() string {
	return fmtx.Yellow.XSPrintf("type=%s source=%s prefix=%s marshal=%s", c.Type, c.Source, c.Prefix, c.Marshal)
}

func (c *Cache) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (c *Cache) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
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
	return nil
}

func Default() *Cache {
	return &Cache{
		Source:  constx.DefaultKey,
		Prefix:  "cache",
		Marshal: "json",
	}
}

func (MultiCache) ID() string {
	return "multi-cache"
}

func (MultiCache) Format() string {
	return ""
}

func (MultiCache) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (m MultiCache) Execute() error {
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
	}
	return nil
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
