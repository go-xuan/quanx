package mongox

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/fmtx"
)

func NewConfigurator(conf *Config) configx.Configurator {
	return conf
}

type Config struct {
	URI             string `json:"uri" yaml:"uri"`                         // 连接uri
	AuthMechanism   string `json:"authMechanism" yaml:"authMechanism"`     // 认证加密方式
	AuthSource      string `json:"authSource" yaml:"authSource"`           // 认证数据库
	Username        string `json:"username" yaml:"username"`               // 用户名
	Password        string `json:"password" yaml:"password"`               // 密码
	Database        string `json:"database" yaml:"database"`               // 数据库名
	MaxPoolSize     uint64 `json:"maxPoolSize" yaml:"maxPoolSize"`         // 连接池最大连接数
	MinPoolSize     uint64 `json:"minPoolSize" yaml:"minPoolSize"`         // 连接池最小连接数
	MaxConnIdleTime int64  `json:"maxConnIdleTime" yaml:"maxConnIdleTime"` // 连接池保持空闲连接的最长时间
}

func (c *Config) Format() string {
	return fmtx.Yellow.XSPrintf("uri=%s database=%s", c.URI, c.Database)
}

func (*Config) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "mongo.yaml",
		NacosDataId: "mongo.yaml",
		Listen:      false,
	}
}

func (c *Config) Execute() error {
	if client, err := c.NewClient(); err != nil {
		log.Error("mongo connect failed: ", c.Format())
		return errorx.Wrap(err, "mongo init client error")
	} else {
		_handler = &Handler{config: c, client: client}
		log.Info("mongo connect successfully: ", c.Format())
		return nil
	}
}

func (c *Config) NewClient() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
	defer cancel()

	// 设置连接选项
	opts := options.Client().ApplyURI(c.URI)
	if c.Username != "" && c.Password != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
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
