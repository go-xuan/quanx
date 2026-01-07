package cachex

import (
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// Config 缓存配置
type Config struct {
	Source   string `json:"source" yaml:"source" default:"default"`     // 缓存源名称
	Client   string `json:"client" yaml:"client" default:"local"`       // 客户端选型（redis/local）
	Enable   bool   `json:"enable" yaml:"enable"`                       // 数据源启用
	Address  string `json:"address" yaml:"address" default:"localhost"` // 主机
	Username string `json:"username" yaml:"username"`                   // 用户名
	Password string `json:"password" yaml:"password"`                   // 密码
	Database int    `json:"database" yaml:"database"`                   // 数据库，默认0
	Prefix   string `json:"prefix" yaml:"prefix"`                       // 缓存key前缀
	Mode     int    `json:"mode" yaml:"mode"`                           // redis模式（0-单机/1-集群/3-哨兵，默认单机模式）
	Master   string `json:"master" yaml:"master"`                       // redis哨兵模式主服务器名称
	PoolSize int    `json:"poolSize" yaml:"poolSize"`                   // redis连接池大小
	Marshal  string `json:"marshal" yaml:"marshal" default:"json"`      // 序列化方式（json/yaml）
}

// Copy 复制配置
func (c *Config) Copy() *Config {
	return &Config{
		Source:   c.Source,
		Client:   c.Client,
		Enable:   c.Enable,
		Address:  c.Address,
		Username: c.Username,
		Password: c.Password,
		Database: c.Database,
		Mode:     c.Mode,
		Prefix:   c.Prefix,
		Marshal:  c.Marshal,
		Master:   c.Master,
	}
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["client"] = c.Client
	fields["address"] = c.Address
	fields["database"] = c.Database
	fields["mode"] = c.Mode
	fields["prefix"] = c.Prefix
	fields["marshal"] = c.Marshal
	return fields
}

func (c *Config) Valid() bool {
	return c.Source != "" && c.Address != ""
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("cache.yaml"),
		configx.NewFileReader("cache.yaml"),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		logger := log.WithFields(c.LogFields())
		client, err := c.NewClient()
		if err != nil {
			logger.WithError(err).Error("cache client init failed")
			return errorx.Wrap(err, "cache client init failed")
		}
		AddClient(client)
		logger.Info("cache client init success")
	}
	return nil
}

// NewClient 根据缓存配置初始化缓存客户端
func (c *Config) NewClient() (Client, error) {
	// 矫正配置
	c.correct()
	// 创建缓存客户端
	client, err := NewClient(c)
	if err != nil {
		return nil, errorx.Wrap(err, "new cache client failed")
	}
	return client, nil
}

// Prefix 矫正配置
func (c *Config) correct() {
	if c.Prefix != "" {
		c.Prefix = stringx.AddSuffix(c.Prefix, ":")
	}
}

// GetKey 获取缓存key
func (c *Config) GetKey(key string) string {
	if c.Prefix != "" {
		return c.Prefix + key
	}
	return key
}

// GetKeys 获取缓存keys
func (c *Config) GetKeys(keys []string) []string {
	if len(keys) > 0 && c.Prefix != "" {
		var newKeys []string
		prefix := stringx.AddSuffix(c.Prefix, ":")
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
