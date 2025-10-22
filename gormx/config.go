package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

// 数据库类型
const (
	MYSQL    = "mysql"    // mysql
	POSTGRES = "postgres" // postgres
	PGSQL    = "pgsql"    // pgsql
)

type Config struct {
	Source          string `json:"source" yaml:"source" default:"default"`              // 数据源名称
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

// LogEntry 日志打印实体类
func (c *Config) LogEntry() *log.Entry {
	return log.WithFields(log.Fields{
		"source":   c.Source,
		"type":     c.Type,
		"host":     c.Host,
		"port":     c.Port,
		"database": c.Database,
	})
}

func (c *Config) Valid() bool {
	return c.Host != "" && c.Port != 0
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("database.yaml"),
		configx.NewFileReader("database.yaml"),
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if client, err := c.NewClient(); err != nil {
			c.LogEntry().WithField("error", err.Error()).Error("database init failed")
			return errorx.Wrap(err, "new gorm client error")
		} else {
			AddClient(client)
			c.LogEntry().Info("database init success")
		}
	}
	return nil
}

func (c *Config) GetLogger() logger.Interface {
	l := DefaultLogger()
	if c.LogLevel != "" {
		l.LogLevel = LogLevel(c.LogLevel)
	}
	if c.SlowThreshold > 0 {
		l.SlowThreshold = time.Duration(c.SlowThreshold) * time.Millisecond
	}
	return l
}

// NewGormDB 创建数据库连接
func (c *Config) NewGormDB() (*gorm.DB, error) {
	var dial gorm.Dialector
	switch strings.ToLower(c.Type) {
	case MYSQL:
		dial = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local",
			c.Username, c.Password, c.Host, c.Port, c.Database))
	case POSTGRES, PGSQL:
		dial = postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			c.Host, c.Port, c.Username, c.Password, c.Database))
	default:
		return nil, errorx.Errorf("database type only support : %v", []string{MYSQL, POSTGRES, PGSQL})
	}
	gormDB, err := gorm.Open(dial, &gorm.Config{
		Logger: c.GetLogger(),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, errorx.Wrap(err, "gorm open failed")
	}
	var sqlDB *sql.DB
	if sqlDB, err = gormDB.DB(); err != nil {
		return nil, errorx.Wrap(err, "get sql db failed")
	}
	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)
	return gormDB, nil
}

func (c *Config) NewClient() (*Client, error) {
	if db, err := c.NewGormDB(); err != nil {
		return nil, errorx.Wrap(err, "new gorm db error")
	} else {
		return &Client{
			config: c,
			db:     db,
		}, nil
	}
}

// CommentTableSql 生成表备注
func (c *Config) CommentTableSql(table, comment string) string {
	switch strings.ToLower(c.Type) {
	case MYSQL:
		return "alter table " + table + " comment = '" + comment + "'"
	case POSTGRES, PGSQL:
		return "comment on table " + table + " is '" + comment + "'"
	}
	return ""
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
			return errorx.Wrap(err, "gorm config execute error")
		}
	}
	if !Initialized() {
		log.Error("database not initialized because no enabled source")
	}
	return nil
}
