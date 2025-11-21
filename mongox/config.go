package mongox

import (
	"context"
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

type Config struct {
	Source          string `json:"source" yaml:"source" default:"default"`                   // 数据源名称
	Enable          bool   `json:"enable" yaml:"enable"`                                     // 数据源启用
	URI             string `json:"uri" yaml:"uri"`                                           // 连接uri
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

func (c *Config) Copy() *Config {
	return &Config{
		Source:          c.Source,
		Enable:          c.Enable,
		URI:             c.URI,
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

// LogEntry 日志打印实体类
func (c *Config) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"source":   c.Source,
		"uri":      c.URI,
		"database": c.Database,
		"debug":    c.Debug,
	})
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("mongo.yaml"),
		configx.NewFileReader("mongo.yaml"),
	}
}

func (c *Config) Valid() bool {
	return c.URI != ""
}

func (c *Config) Execute() error {
	if c.Enable {
		if client, err := c.NewClient(); err != nil {
			c.LogEntry().WithError(err).Error("mongo init failed")
			return errorx.Wrap(err, "mongo init client error")
		} else {
			c.LogEntry().Info("mongo init success")
			AddClient(c, client)
		}
	}
	return nil
}

func (c *Config) NewClient() (*mongo.Client, error) {
	ctx := context.TODO()
	if c.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(c.Timeout)*time.Millisecond)
		defer cancel()
	}
	// 设置连接选项
	opts := options.Client().ApplyURI(c.URI)
	if c.Username != "" && c.Password != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: c.AuthMechanism,
			AuthSource:    c.AuthSource,
			Username:      c.Username,
			Password:      c.Password,
		})
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
	opts.SetReadPreference(readpref.PrimaryPreferred())

	if c.Debug {
		opts.SetMonitor(&event.CommandMonitor{
			Started: func(ctx context.Context, event *event.CommandStartedEvent) {
				if name := event.CommandName; name != "ping" {
					log.WithField("command_name", event.CommandName).
						WithField("command", event.Command.String()).
						Info("mongo command start")
				}
			},
			Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
				if name := event.CommandName; name != "ping" {
					log.WithField("command_name", event.CommandName).
						WithField("duration", event.Duration.String()).
						Info("mongo command success")
				}
			},
			Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
				if name := event.CommandName; name != "ping" {
					log.WithField("command_name", event.CommandName).
						WithField("duration", event.Duration.String()).
						Info("mongo command failed")
				}
			},
		})
	}

	// 建立连接
	if client, err := mongo.Connect(ctx, opts); err != nil {
		return nil, errorx.Wrap(err, "mongo connect failed")
	} else if err = client.Ping(ctx, readpref.PrimaryPreferred()); err != nil {
		return nil, errorx.Wrap(err, "mongo ping failed")
	} else {
		return client, nil
	}
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
