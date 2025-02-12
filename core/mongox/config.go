package mongox

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/anyx"
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
	MaxConnIdleTime int64  `json:"maxConnIdleTime" yaml:"maxConnIdleTime"`                   // 连接池保持空闲连接的最长时间
}

func (c *Config) Format() string {
	return fmt.Sprintf("uri=%s database=%s", c.URI, c.Database)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if client, err := c.NewClient(); err != nil {
			log.Error("mongo connect failed: ", c.Format())
			return errorx.Wrap(err, "mongo init client error")
		} else {
			log.Info("mongo connect success: ", c.Format())
			if _handler == nil {
				_handler = &Handler{
					multi: false, config: c, client: client,
					configs: make(map[string]*Config),
					clients: make(map[string]*mongo.Client),
				}
			} else {
				_handler.multi = true
				if c.Source == constx.DefaultSource {
					_handler.config = c
					_handler.client = client
				}
			}
			_handler.configs[c.Source] = c
			_handler.clients[c.Source] = client
		}
	}
	return nil
}

func (c *Config) NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

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
		opts.SetMaxConnIdleTime(time.Second * time.Duration(c.MaxConnIdleTime))
	}
	opts.SetReadPreference(readpref.PrimaryPreferred())

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

func (m MultiConfig) Format() string {
	sb := &strings.Builder{}
	sb.WriteString("[")
	for i, config := range m {
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

func (MultiConfig) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

func (m MultiConfig) Execute() error {
	if len(m) == 0 {
		return errorx.New("mongo not connected! cause: mongo.yaml is invalid")
	}
	if _handler == nil {
		_handler = &Handler{
			configs: make(map[string]*Config),
			clients: make(map[string]*mongo.Client),
		}
	}
	_handler.multi = true
	for i, c := range m {
		if c.Enable {
			if err := anyx.SetDefaultValue(c); err != nil {
				return errorx.Wrap(err, "set default value error")
			}
			if client, err := c.NewClient(); err != nil {
				return errorx.Wrap(err, "new mongo client failed")
			} else {
				_handler.clients[c.Source] = client
				_handler.configs[c.Source] = c
				if i == 0 || c.Source == constx.DefaultSource {
					_handler.client = client
					_handler.config = c
				}
			}
		}
	}
	if len(_handler.configs) == 0 {
		log.Error("mongo not connected! cause: mongo.yaml is empty or no enabled source")
	}
	return nil
}
