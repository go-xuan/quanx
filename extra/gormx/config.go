package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/types/anyx"
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
	Debug           bool   `json:"debug" yaml:"debug" default:"false"`                  // 开启debug（打印SQL以及初始化模型建表）
	MaxIdleConns    int    `json:"maxIdleConns" yaml:"maxIdleConns" default:"10"`       // 最大空闲连接
	MaxOpenConns    int    `json:"maxOpenConns" yaml:"maxOpenConns" default:"10"`       // 最大打开连接
	ConnMaxLifetime int    `json:"connMaxLifetime" yaml:"connMaxLifetime" default:"10"` // 连接存活时间(分钟)
}

func (c *Config) Format() string {
	return fmt.Sprintf("source=%s type=%s host=%s port=%v database=%s debug=%v",
		c.Source, c.Type, c.Host, c.Port, c.Database, c.Debug)
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "database.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "database.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) FileReader() configx.Reader {
	return &configx.LocalFileReader{
		Name: "database.yaml",
	}
}

func (c *Config) Execute() error {
	if c.Enable {
		if err := anyx.SetDefaultValue(c); err != nil {
			return errorx.Wrap(err, "set default value error")
		}
		if db, err := c.NewGormDB(); err != nil {
			return errorx.Wrap(err, "new gorm db error")
		} else {
			if _handler == nil {
				_handler = &Handler{
					multi: false, config: c, db: db,
					configs: make(map[string]*Config),
					dbs:     make(map[string]*gorm.DB),
				}
			} else {
				_handler.multi = true
				if c.Source == constx.DefaultSource {
					_handler.config = c
					_handler.db = db
				}
			}
			_handler.configs[c.Source] = c
			_handler.dbs[c.Source] = db
		}
	}
	return nil
}

// NewGormDB 创建数据库连接
func (c *Config) NewGormDB() (*gorm.DB, error) {
	if db, err := c.GetGormDB(); err != nil {
		return nil, errorx.Wrap(err, "new gorm db failed")
	} else {
		var sqlDB *sql.DB
		if sqlDB, err = db.DB(); err != nil {
			return nil, errorx.Wrap(err, "get sql db failed")
		}
		sqlDB.SetMaxIdleConns(c.MaxIdleConns)
		sqlDB.SetMaxOpenConns(c.MaxOpenConns)
		sqlDB.SetConnMaxLifetime(time.Duration(c.ConnMaxLifetime) * time.Second)

		if c.Debug {
			// 是否打印SQL
			db = db.Debug()
		}
		return db, nil
	}
}

// 数据库类型
const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
	PGSQL    = "pgsql"
)

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

// GetGormDB 根据dsn生成gormDB
func (c *Config) GetGormDB() (*gorm.DB, error) {
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
	if db, err := gorm.Open(dial, &gorm.Config{
		// 模型命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名单数命名
		},
	}); err != nil {
		return nil, errorx.Wrap(err, "gorm open dialector failed")
	} else {
		return db, nil
	}
}

// MultiConfig 数据库多数据源配置
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

func (MultiConfig) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FormNacos:
		return &nacosx.Reader{
			DataId: "database.yaml",
		}
	case configx.FromLocal:
		return &configx.LocalFileReader{
			Name: "database.yaml",
		}
	default:
		return nil
	}
}

func (m MultiConfig) Execute() error {
	if len(m) == 0 {
		return errorx.New("database not connected! cause: database.yaml is invalid")
	}
	if _handler == nil {
		_handler = &Handler{
			dbs:     make(map[string]*gorm.DB),
			configs: make(map[string]*Config),
		}
	}
	_handler.multi = true
	for i, c := range m {
		if c.Enable {
			if err := anyx.SetDefaultValue(c); err != nil {
				return errorx.Wrap(err, "set default value error")
			}
			db, err := c.NewGormDB()
			if err != nil {
				return errorx.Wrap(err, "new gorm.Config failed")
			}
			_handler.dbs[c.Source] = db
			_handler.configs[c.Source] = c
			if i == 0 || c.Source == constx.DefaultSource {
				_handler.db = db
				_handler.config = c
			}
		}
	}
	if len(_handler.configs) == 0 {
		log.Error("database not connected! cause: database.yaml is empty or no enabled source")
	}
	return nil
}
