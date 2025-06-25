package mongox

import (
	"context"
	"fmt"
	configx2 "github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"strings"
	"time"

	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

func (c *Config) Info() string {
	return fmt.Sprintf("uri=%s database=%s", c.URI, c.Database)
}

func (*Config) Reader(from configx2.From) configx2.Reader {
	switch from {
	case configx2.FromNacos:
		return &nacosx.Reader{
			DataId: "mongo.yaml",
		}
	case configx2.FromFile:
		return &configx2.FileReader{
			Name: "mongo.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if client, err := c.NewClient(); err != nil {
			log.Error("mongo connect failed: ", c.Info())
			return errorx.Wrap(err, "mongo init client error")
		} else {
			log.Info("mongo connect success: ", c.Info())
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
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				if startedEvent.CommandName != "ping" {
					log.Info("current mongo command: ", startedEvent.Command)
				}
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				if succeededEvent.CommandName != "ping" {
					log.Infof("mongo command success: %s ==> %v", succeededEvent.CommandName, succeededEvent.Duration)
				}
			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				if failedEvent.CommandName != "ping" {
					log.Errorf("mongo command failed: %s ==> %v", failedEvent.CommandName, failedEvent.Duration)
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

type MultiConfig []*Config

func (list MultiConfig) Info() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range list {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("{")
		sb.WriteString(config.Info())
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}

func (MultiConfig) Reader(from configx2.From) configx2.Reader {
	switch from {
	case configx2.FromNacos:
		return &nacosx.Reader{
			DataId: "mongo.yaml",
		}
	case configx2.FromFile:
		return &configx2.FileReader{
			Name: "mongo.yaml",
		}
	default:
		return nil
	}
}

func (list MultiConfig) Execute() error {
	if len(list) == 0 {
		return errorx.New("mongo not initialized, mongo.yaml is invalid")
	}
	for _, config := range list {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "mongo config execute error")
		}
	}
	if !Initialized() {
		log.Error("mongo not initialized, no enabled source")
	}
	return nil
}
