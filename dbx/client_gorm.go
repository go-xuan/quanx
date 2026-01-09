package dbx

import (
	"database/sql"
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

// NewGormClient 创建gorm客户端
func NewGormClient(config *Config) (*GormClient, error) {
	db, err := NewGormDB(config)
	if err != nil {
		return nil, errorx.Wrap(err, "new gorm client failed")
	}
	return &GormClient{config: config, db: db}, nil
}

// GormClient gorm客户端
type GormClient struct {
	config *Config
	db     *gorm.DB
}

func (c *GormClient) GetClient() *gorm.DB {
	return c.db
}

func (c *GormClient) GetInstance() interface{} {
	return c.db
}

func (c *GormClient) GetConfig() *Config {
	return c.config
}

func (c *GormClient) Close() error {
	logger := log.WithFields(c.config.LogFields())
	if err := CloseGormDB(c.db); err != nil {
		logger.WithError(err).Error("gorm close failed")
		return errorx.Wrap(err, "gorm close failed")
	}
	logger.Info("db client close success")
	return nil
}

func (c *GormClient) Raw(sql string, dest interface{}) error {
	raw := c.GetClient().Raw(sql)
	if err := raw.Error; err != nil {
		return errorx.Wrap(err, "gorm raw failed")
	}
	if dest != nil {
		if err := raw.Scan(dest).Error; err != nil {
			return errorx.Wrap(err, "scan to dest failed")
		}
	}
	return nil
}

func (c *GormClient) Exec(sql string) error {
	if err := c.GetClient().Exec(sql).Error; err != nil {
		return errorx.Wrap(err, "gorm exec failed")
	}
	return nil
}

// NewGormDB 创建gorm数据库连接
func NewGormDB(config *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch config.Dialect {
	case MYSQL:
		dialector = mysql.Open(config.GetDSN())
	case POSTGRES, PGSQL:
		dialector = postgres.Open(config.GetDSN())
	default:
		return nil, errorx.New("unsupported dialect: " + config.Dialect)
	}
	db, err := gorm.Open(dialector, &gorm.Config{
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
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifetime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(config.MaxIdleTime) * time.Second)
	return db, nil
}

// CloseGormDB 关闭gorm数据库连接
func CloseGormDB(db *gorm.DB) error {
	if d, err := db.DB(); err != nil {
		return errorx.Wrap(err, "get sql db failed")
	} else if err = d.Close(); err != nil {
		return errorx.Wrap(err, "close sql db failed")
	}
	return nil
}
