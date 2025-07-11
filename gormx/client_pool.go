package gormx

import (
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"
)

var pool *typex.Enum[string, *Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil
}

func this() *typex.Enum[string, *Client] {
	if !Initialized() {
		panic("gorm client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, db *gorm.DB) {
	if config == nil || db == nil {
		return
	}
	client := &Client{config, db}
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		this().Add("default", client)
	}
	this().Add(config.Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get("default")
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

// CopyDatabase 复制数据库
func CopyDatabase(source, target, database string) error {
	if source != "" && target != "" {
		if sourceClient := this().Get(source); sourceClient != nil {
			if targetClient, err := sourceClient.Copy(target, database); err == nil {
				this().Add(target, targetClient)
			} else {
				return errorx.Wrap(err, "copy database client failed")
			}
		} else {
			return errorx.New("database source client config not exist")
		}
	}
	return nil
}
