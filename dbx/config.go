package dbx

import (
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
)

// 数据库类型
const (
	MYSQL    = "mysql"    // mysql
	POSTGRES = "postgres" // postgres
	PGSQL    = "pgsql"    // pgsql
)

type Config struct {
	Source          string `json:"source" yaml:"source" default:"default"`              // 数据源名称
	Client          string `json:"client" yaml:"client" default:"gorm"`                 // 客户端选型
	Enable          bool   `json:"enable" yaml:"enable"`                                // 数据源启用
	Type            string `json:"type" yaml:"type"`                                    // 数据库类型
	Host            string `json:"host" yaml:"host" default:"localhost"`                // 数据库Host
	Port            int    `json:"port" yaml:"port"`                                    // 数据库端口
	Username        string `json:"username" yaml:"username"`                            // 用户名
	Password        string `json:"password" yaml:"password"`                            // 密码
	Database        string `json:"database" yaml:"database"`                            // 数据库名
	Schema          string `json:"schema" yaml:"schema"`                                // schema模式名
	MaxIdleConns    int    `json:"maxIdleConns" yaml:"maxIdleConns" default:"10"`       // 最大空闲连接
	MaxOpenConns    int    `json:"maxOpenConns" yaml:"maxOpenConns" default:"100"`      // 最大打开连接
	ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime" default:"10"` // 连接存活时间(秒)
	LogLevel        string `json:"logLevel" yaml:"logLevel" default:"warn"`             // 日志级别
	SlowThreshold   int    `json:"slowThreshold" yaml:"slowThreshold" default:"200"`    // 慢查询阈值(毫秒)
}

func (c *Config) Copy() *Config {
	return &Config{
		Source:          c.Source,
		Client:          c.Client,
		Enable:          c.Enable,
		Type:            c.Type,
		Host:            c.Host,
		Port:            c.Port,
		Username:        c.Username,
		Password:        c.Password,
		Database:        c.Database,
		Schema:          c.Schema,
		MaxIdleConns:    c.MaxIdleConns,
		MaxOpenConns:    c.MaxOpenConns,
		ConnMaxLifetime: c.ConnMaxLifetime,
		LogLevel:        c.LogLevel,
		SlowThreshold:   c.SlowThreshold,
	}
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["client"] = c.Client
	fields["type"] = c.Type
	fields["host"] = c.Host
	fields["port"] = c.Port
	fields["database"] = c.Database
	return fields
}

func (c *Config) Valid() bool {
	return c.Source != "" && c.Host != "" && c.Port > 0
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("database.yaml"),
		configx.NewFileReader("database.yaml"),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		logger_ := log.WithFields(c.LogFields())
		client, err := c.NewClient()
		if err != nil {
			logger_.WithError(err).Error("database client init failed")
			return errorx.Wrap(err, "init database client failed")
		}
		AddClient(client)
		logger_.Info("database client init success")
	}
	return nil
}

func (c *Config) NewClient() (Client, error) {
	client, err := NewClient(c)
	if err != nil {
		return nil, errorx.Wrap(err, "new database client failed")
	}
	return client, nil
}

type Configs []*Config

func (s Configs) Valid() bool {
	return len(s) > 0
}

func (s Configs) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("database.yaml"),
		configx.NewFileReader("database.yaml"),
	}
}

func (s Configs) Execute() error {
	for _, config := range s {
		if err := config.Execute(); err != nil {
			return errorx.Wrap(err, "database config execute error")
		}
	}
	if !Initialized() {
		log.Error("database not initialized because no enabled source")
	}
	return nil
}
