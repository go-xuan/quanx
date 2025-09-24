package cachex

import (
	"time"

	"github.com/go-xuan/quanx/logx"
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
	Type    string `json:"type" yaml:"type" default:"local"`       // 缓存类型（local/redis）
	Source  string `json:"source" yaml:"source" default:"default"` // 缓存存储数据源名称
	Prefix  string `json:"prefix" yaml:"prefix" default:"local"`   // 缓存key前缀前缀
	Marshal string `json:"marshal" yaml:"marshal" default:"json"`  // 序列化方案
}

func (c *Config) LogFields() log.Fields {
	return log.Fields{
		"type":    c.Type,
		"source":  c.Source,
		"prefix":  c.Prefix,
		"marshal": c.Marshal,
	}
}

func (c *Config) Valid() bool {
	return c.Type != "" && c.Source != ""
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
		logx.WithEntity(c).WithField("error", err.Error()).Error("cache init failed")
		return errorx.Wrap(err, "new cache error")
	} else {
		logx.WithEntity(c).Info("cache init success")
		AddClient(c, client)
	}
	return nil
}

// NewClient 根据缓存配置初始化缓存客户端
func (c *Config) NewClient() (Client, error) {
	switch c.Type {
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
