package dbx

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// 数据库类型
const (
	MYSQL    = "mysql"    // mysql
	POSTGRES = "postgres" // postgres
	PGSQL    = "pgsql"    // pgsql
)

// Config 数据库配置 dialect
type Config struct {
	Source        string            `json:"source" yaml:"source" default:"default"`           // 数据源名称
	Builder       string            `json:"builder" yaml:"builder" default:"gorm"`            // 客户端选型
	Enable        bool              `json:"enable" yaml:"enable"`                             // 数据源启用
	Dialect       string            `json:"dialect" yaml:"dialect"`                           // 数据库方言
	Dsn           string            `json:"dsn" yaml:"dsn"`                                   // DSN连接字符串
	Host          string            `json:"host" yaml:"host" default:"localhost"`             // 数据库Host
	Port          int               `json:"port" yaml:"port"`                                 // 数据库端口
	Username      string            `json:"username" yaml:"username"`                         // 用户名
	Password      string            `json:"password" yaml:"password"`                         // 密码
	Database      string            `json:"database" yaml:"database"`                         // 数据库名
	Schema        string            `json:"schema" yaml:"schema"`                             // schema模式名
	Options       map[string]string `json:"options" yaml:"options"`                           // 连接选项，可覆盖默认值
	MaxOpenConns  int               `json:"maxOpenConns" yaml:"maxOpenConns" default:"100"`   // 最大打开连接
	MaxIdleConns  int               `json:"maxIdleConns" yaml:"maxIdleConns" default:"10"`    // 最大空闲连接
	MaxLifetime   int               `json:"maxLifetime" yaml:"maxLifetime" default:"10"`      // 连接存活时间(秒)
	MaxIdleTime   int               `json:"maxIdleTime" yaml:"maxIdleTime" default:"10"`      // 连接空闲时间(毫秒)
	LogLevel      string            `json:"logLevel" yaml:"logLevel" default:"warn"`          // 日志级别
	SlowThreshold int               `json:"slowThreshold" yaml:"slowThreshold" default:"200"` // 慢查询阈值(毫秒)
}

// GetDSN 获取DSN
func (c *Config) GetDSN() string {
	if c.Dsn == "" {
		dsn := strings.Builder{}
		switch c.Dialect {
		case MYSQL:
			options := make(map[string]string)
			options["clientFoundRows"] = "false"
			options["timeout"] = "10s"
			options["charset"] = "utf8"
			options["collation"] = "utf8_general_ci"
			options["parseTime"] = "true"
			options["loc"] = "Asia/Shanghai"
			for k, v := range c.Options {
				options[k] = v
			}
			dsn.WriteString(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database))
			if len(options) > 0 {
				dsn.WriteString("?")
				conn := ""
				for k, v := range options {
					dsn.WriteString(conn)
					dsn.WriteString(k)
					dsn.WriteString("=")
					dsn.WriteString(url.QueryEscape(v))
					conn = "&"
				}
			}
		case POSTGRES, PGSQL:
			dsn.WriteString(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
				c.Host, c.Port, c.Username, c.Password, c.Database))
			options := make(map[string]string)
			options["sslmode"] = "disable"
			options["TimeZone"] = "Asia/Shanghai"
			for k, v := range c.Options {
				options[k] = v
			}
			for k, v := range options {
				dsn.WriteString(fmt.Sprintf(" %s=%s", k, v))
			}
		}
		c.Dsn = dsn.String()
	}
	return c.Dsn
}

// Copy 复制配置
func (c *Config) Copy() *Config {
	return &Config{
		Source:        c.Source,
		Builder:       c.Builder,
		Enable:        c.Enable,
		Dialect:       c.Dialect,
		Host:          c.Host,
		Port:          c.Port,
		Username:      c.Username,
		Password:      c.Password,
		Database:      c.Database,
		Schema:        c.Schema,
		MaxOpenConns:  c.MaxOpenConns,
		MaxIdleConns:  c.MaxIdleConns,
		MaxLifetime:   c.MaxLifetime,
		MaxIdleTime:   c.MaxIdleTime,
		LogLevel:      c.LogLevel,
		SlowThreshold: c.SlowThreshold,
	}
}

// LogFields 日志字段
func (c *Config) LogFields() map[string]interface{} {
	fields := make(map[string]interface{})
	fields["source"] = c.Source
	fields["builder"] = c.Builder
	fields["dialect"] = c.Dialect
	fields["host"] = c.Host
	fields["port"] = c.Port
	fields["database"] = c.Database
	return fields
}

func (c *Config) Valid() bool {
	return c.Dsn != "" || (c.Host != "" && c.Port > 0)
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
		client, err := NewClient(c)
		if err != nil {
			logger_.WithError(err).Error("init database client failed")
			return errorx.Wrap(err, "init database client failed")
		}
		AddClient(c.Source, client)
		logger_.Info("init database client success")
	}
	return nil
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
