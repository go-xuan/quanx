package gormx

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
)

var pool *enumx.Enum[string, *Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

func this() *enumx.Enum[string, *Client] {
	if pool == nil {
		panic("gorm client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, db *gorm.DB) {
	if pool == nil {
		pool = enumx.NewStringEnum[*Client]()
	}
	pool.Add(config.Source, &Client{config, db})
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	return this().Get(stringx.Default(constx.DefaultSource, source...))
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *gorm.DB {
	return GetClient(source...).Instance()
}

// Sources 所有数据源
func Sources() []string {
	return this().Keys()
}

// Close 关闭数据库连接
func Close(source ...string) error {
	if db, err := GetInstance(source...).DB(); err != nil {
		return errorx.Wrap(err, "get db failed")
	} else if err = db.Close(); err != nil {
		return errorx.Wrap(err, "close db failed")
	}
	return nil
}

// InitTable 初始化表结构以及表数据
func InitTable(source string, tablers ...interface{}) error {
	return initTablers(source, tablers...)
}

// CopyDatabase 复制数据库
func CopyDatabase(source, target, database string) {
	if source != "" && target != "" {
		logger := log.WithField("source", source).
			WithField("target", target).
			WithField("database", database)
		if client := this().Get(source); client != nil {
			config := client.Config().Copy()
			config.Source = target
			config.Database = database
			if db, err := config.NewGormDB(); err == nil {
				AddClient(config, db)
			} else {
				logger.WithError(err).Error("copy database connection failed")
			}
		} else {
			logger.Error("database config not found")
		}
	}
}
