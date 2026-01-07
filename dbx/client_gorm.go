package dbx

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
	"gorm.io/gorm/schema"
)

// GormClientBuilder gorm客户端构建器
func GormClientBuilder(config *Config) (Client, error) {
	return NewGormClient(config)
}

// NewGormClient gorm客户端构建器
func NewGormClient(config *Config) (*GormClient, error) {
	db, err := NewGormDB(config)
	if err != nil {
		return nil, errorx.Wrap(err, "new gorm client failed")
	}
	return &GormClient{
		config: config,
		db:     db,
	}, nil
}

// GormClient gorm客户端
type GormClient struct {
	config *Config
	db     *gorm.DB
}

func (c *GormClient) GetConfig() *Config {
	return c.config
}

func (c *GormClient) GetInstance() interface{} {
	return c.db
}

func (c *GormClient) Copy(target, database string) (Client, error) {
	config := c.config.Copy()
	config.Source = target
	config.Database = database

	logger := log.WithFields(config.LogFields())
	client, err := config.NewClient()
	if err != nil {
		logger.WithError(err).Error("gorm client copy failed")
		return nil, errorx.Wrap(err, "gorm client copy failed")
	}
	logger.Info("gorm client copy success")
	return client, nil
}

func (c *GormClient) Close() error {
	if err := CloseGormDB(c.db); err != nil {
		return errorx.Wrap(err, "close gorm db failed")
	}
	return nil
}

// CloseGormDB 关闭gorm数据库连接
func CloseGormDB(db *gorm.DB) error {
	if d, err := db.DB(); err != nil {
		return errorx.Wrap(err, "get db failed")
	} else if err = d.Close(); err != nil {
		return errorx.Wrap(err, "close db failed")
	}
	return nil
}

// NewGormDB 创建数据库连接
func NewGormDB(config *Config) (*gorm.DB, error) {
	var dial gorm.Dialector
	switch strings.ToLower(config.Type) {
	case MYSQL:
		dial = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?clientFoundRows=false&parseTime=true&timeout=1800s&charset=utf8&collation=utf8_general_ci&loc=Local",
			config.Username, config.Password, config.Host, config.Port, config.Database))
	case POSTGRES, PGSQL:
		dial = postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.Host, config.Port, config.Username, config.Password, config.Database))
	default:
		return nil, errorx.Sprintf("database type only support: %v", []string{MYSQL, POSTGRES, PGSQL})
	}
	db, err := gorm.Open(dial, &gorm.Config{
		Logger: NewGormLogger(config.LogLevel, time.Duration(config.SlowThreshold)*time.Millisecond),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, errorx.Wrap(err, "gorm open failed")
	}
	var sqlDB *sql.DB
	if sqlDB, err = db.DB(); err != nil {
		return nil, errorx.Wrap(err, "get sql db failed")
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	return db, nil
}
