package cachex

import (
	"fmt"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/extra/redisx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

const (
	CacheTypeRedis = "redis"
	CacheTypeLocal = "local"
)

type Config struct {
	Type    string `json:"type" yaml:"type" default:"redis"`         // 缓存类型（local/redis）
	Source  string `json:"source" yaml:"source" default:"default"`   // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"default"`   // 缓存KEY前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"msgpack"` // 序列化方案
}

func (c *Config) Format() string {
	return fmt.Sprintf("type=%s source=%s prefix=%s marshal=%s",
		c.Type, c.Source, c.Prefix, c.Marshal)
}

func (c *Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "cache.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "cache.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if err := anyx.SetDefaultValue(c); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	if client, err := c.NewClient(); err != nil {
		log.Error("cache init failed: ", c.Format())
		return errorx.Wrap(err, "new client error")
	} else {
		log.Info("cache init success: ", c.Format())
		AddClient(c, client)
	}
	return nil
}

// NewClient 根据缓存配置初始化缓存客户端
func (c *Config) NewClient() (Client, error) {
	switch c.Type {
	case CacheTypeRedis:
		return &RedisClient{
			config:  c,
			client:  redisx.GetClient(c.Source),
			marshal: marshalx.Apply(c.Marshal),
		}, nil
	case CacheTypeLocal:
		return &LocalClient{
			config:  c,
			client:  cache.New(time.Duration(-1), time.Duration(-1)),
			marshal: marshalx.Apply(c.Marshal),
		}, nil
	default:
		return nil, errorx.New("not support type: " + c.Type)
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

// MultiConfig 多缓存配置
type MultiConfig []*Config

func (list MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range list {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("{")
		sb.WriteString(config.Format())
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}

func (MultiConfig) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "cache.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "cache.yaml",
		}
	default:
		return nil
	}
}

func (list MultiConfig) Execute() error {
	if len(list) == 0 {
		return errorx.New("cache client not init! cause: cache.yaml is invalid")
	}
	for _, config := range list {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "cache config execute error")
		}
	}
	if len(_handler.clients) == 0 {
		log.Error("cache client not init! cause: no enabled source")
	}
	return nil
}
