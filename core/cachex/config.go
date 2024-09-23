package cachex

import (
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/redisx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/fmtx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

const (
	CacheTypeRedis = "redis"
	CacheTypeLocal = "local"
)

// MultiConfig 多缓存配置
type MultiConfig []*Config

type Config struct {
	Type    string `json:"type" yaml:"type" default:"redis"`         // 缓存类型（local/redis）
	Source  string `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
}

func (c *Config) ID() string {
	return "cache"
}

func (c *Config) Format() string {
	return fmtx.Yellow.XSPrintf("type=%s source=%s prefix=%s marshal=%s",
		c.Type, c.Source, c.Prefix, c.Marshal)
}

func (c *Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	var client = c.InitClient()
	if _handler == nil {
		_handler = &Handler{
			multi:     false,
			client:    client,
			clientMap: make(map[string]Client),
		}
	} else {
		_handler.multi = true
	}
	_handler.clientMap[c.Source] = client
	return nil
}

func Default() *Config {
	return &Config{
		Source:  constx.DefaultSource,
		Prefix:  "cache",
		Marshal: "json",
	}
}

func (MultiConfig) ID() string {
	return "multi-cache"
}

func (MultiConfig) Format() string {
	return ""
}

func (MultiConfig) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "cache.yaml",
		NacosDataId: "cache.yaml",
		Listen:      false,
	}
}

func (m MultiConfig) Execute() error {
	if _handler == nil {
		_handler = &Handler{
			multi:     true,
			clientMap: make(map[string]Client),
		}
	} else {
		_handler.multi = true
	}
	multi := anyx.IfZero(m, MultiConfig{Default()})
	for i, c := range multi {
		var client = c.InitClient()
		_handler.clientMap[c.Source] = client
		if i == 0 || c.Source == constx.DefaultSource {
			_handler.client = client
		}
	}
	return nil
}

// InitClient 根据缓存配置初始化缓存客户端
func (c *Config) InitClient() Client {
	switch c.Type {
	case CacheTypeRedis:
		return &RedisClient{
			config:  c,
			client:  redisx.Client(c.Source),
			marshal: marshalx.NewCase(c.Marshal),
		}
	case CacheTypeLocal:
		return &LocalClient{
			config:  c,
			client:  cache.New(time.Duration(-1), time.Duration(-1)),
			convert: marshalx.NewCase(c.Marshal),
		}
	default:
		log.Error("cache client not support type: ", c.Type)
		return nil
	}
}

func (c *Config) GetKey(key string) string {
	if c.Prefix != "" {
		return stringx.AddSuffix(c.Prefix, ":") + key
	} else {
		return key
	}
}

func (c *Config) GetKeys(keys []string) []string {
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
