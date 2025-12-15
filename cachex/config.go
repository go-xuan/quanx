package cachex

import (
	"time"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/go-xuan/utilx/stringx"
	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
)

const (
	REDIS = "redis"
	LOCAL = "local"
)

type Config struct {
	Source  string `json:"source" yaml:"source" default:"default"` // 缓存源名称
	Storage string `json:"storage" yaml:"storage" default:"redis"` // 缓存存储类型（redis/local）
	Prefix  string `json:"prefix" yaml:"prefix" default:"local"`   // 缓存key前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"json"`  // 序列化方式（json/yaml）
}

func (c *Config) LogFields() log.Fields {
	return log.Fields{
		"source":  c.Source,
		"storage": c.Storage,
		"prefix":  c.Prefix,
		"marshal": c.Marshal,
	}
}

func (c *Config) Valid() bool {
	return c.Storage != "" && c.Source != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("cache.yaml"),
		configx.NewFileReader("cache.yaml"),
		configx.NewTagReader(),
	}
}

func (c *Config) Execute() error {
	if client, err := c.NewClient(); err != nil {
		log.WithFields(c.LogFields()).WithError(err).Error("cache client init failed")
		return errorx.Wrap(err, "new cache client error")
	} else {
		log.WithFields(c.LogFields()).Info("cache client init success")
		AddClient(c.Source, client)
	}
	return nil
}

// NewClient 根据缓存配置初始化缓存客户端
func (c *Config) NewClient() (Client, error) {
	switch c.Storage {
	case REDIS:
		return &RedisClient{
			config:  c,
			client:  redisx.GetInstance(c.Source),
			marshal: marshalx.Apply(c.Marshal),
		}, nil
	case LOCAL:
		return &LocalClient{
			config:  c,
			cache:   cache.New(time.Duration(-1), time.Duration(-1)),
			marshal: marshalx.Apply(c.Marshal),
		}, nil
	default:
		return nil, errorx.Sprintf("not support cache storage: %s", c.Storage)
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

type Configs []*Config

func (s Configs) Valid() bool {
	return len(s) > 0
}

func (s Configs) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("cache.yaml"),
		configx.NewFileReader("cache.yaml"),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "cache config execute error")
		}
	}
	if !Initialized() {
		log.Error("cache not initialized because no enabled source")
	}
	return nil
}
