package mongox

import (
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

type Config struct {
	Source          string `json:"source" yaml:"source" default:"default"`                   // 数据源名称
	Enable          bool   `json:"enable" yaml:"enable"`                                     // 数据源启用
	Uri             string `json:"uri" yaml:"uri"`                                           // 连接uri
	AuthMechanism   string `json:"authMechanism" yaml:"authMechanism" default:"SCRAM-SHA-1"` // 认证加密方式
	AuthSource      string `json:"authSource" yaml:"authSource"`                             // 认证数据库
	Username        string `json:"username" yaml:"username"`                                 // 用户名
	Password        string `json:"password" yaml:"password"`                                 // 密码
	Database        string `json:"database" yaml:"database"`                                 // 数据库名
	MaxPoolSize     uint64 `json:"maxPoolSize" yaml:"maxPoolSize"`                           // 连接池最大连接数
	MinPoolSize     uint64 `json:"minPoolSize" yaml:"minPoolSize"`                           // 连接池最小连接数
	MaxConnIdleTime uint64 `json:"maxConnIdleTime" yaml:"maxConnIdleTime"`                   // 连接池保持空闲连接的最长时间
	Timeout         uint64 `json:"timeout" yaml:"timeout"`                                   // 超时时间
	Debug           bool   `json:"debug" yaml:"debug"`                                       // debug模式（日志打印）
}

// ClientOptions 连接选项
func (c *Config) ClientOptions() *options.ClientOptions {
	opts := options.Client().
		ApplyURI(c.Uri).                               // 设置uri
		SetReadPreference(readpref.PrimaryPreferred()) // 设置读取偏好为主节点优先

	// 设置认证信息
	if c.Username != "" && c.Password != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: c.AuthMechanism,
			AuthSource:    c.AuthSource,
			Username:      c.Username,
			Password:      c.Password,
		})
	}
	if c.Timeout > 0 {
		opts.SetTimeout(time.Duration(c.Timeout) * time.Millisecond)
	}
	if c.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(c.MaxPoolSize)
	}
	if c.MinPoolSize > 0 {
		opts.SetMinPoolSize(c.MinPoolSize)
	}
	if c.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(time.Duration(c.MaxConnIdleTime) * time.Millisecond)
	}
	if c.Debug {
		opts.SetMonitor(DebugCommandMonitor())
	}
	return opts
}

func (c *Config) Copy() *Config {
	return &Config{
		Source:          c.Source,
		Enable:          c.Enable,
		Uri:             c.Uri,
		AuthMechanism:   c.AuthMechanism,
		AuthSource:      c.AuthSource,
		Username:        c.Username,
		Password:        c.Password,
		Database:        c.Database,
		MaxPoolSize:     c.MaxPoolSize,
		MinPoolSize:     c.MinPoolSize,
		MaxConnIdleTime: c.MaxConnIdleTime,
		Timeout:         c.Timeout,
		Debug:           c.Debug,
	}
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["uri"] = c.Uri
	fields["database"] = c.Database
	fields["debug"] = c.Debug
	return fields
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("mongo.yaml"),
		configx.NewFileReader("mongo.yaml"),
	}
}

func (c *Config) Valid() bool {
	return c.Uri != ""
}

func (c *Config) Execute() error {
	if c.Enable {
		logger := log.WithFields(c.LogFields())
		client, err := NewClient(c)
		if err != nil {
			logger.WithError(err).Error("init mongo client failed")
			return errorx.Wrap(err, "init mongo client failed")
		}
		logger.Info("init mongo client success")
		AddClient(c.Source, client)
	}
	return nil
}

type Configs []*Config

func (s Configs) Valid() bool {
	return len(s) > 0
}

func (s Configs) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("mongo.yaml"),
		configx.NewFileReader("mongo.yaml"),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "mongo config execute error")
		}
	}
	if !Initialized() {
		log.Error("mongo not initialized because no enabled source")
	}
	return nil
}
