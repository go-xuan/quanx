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
		panic("gorm client pool not initialized, please check the relevant config")
	}
	return pool
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

// AddClient 添加客户端
func AddClient(client *Client) {
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		this().Add("default", client)
	}
	this().Add(client.GetConfig().Source, client)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *gorm.DB {
	return GetClient(source...).GetInstance()
}

// Sources 所有数据源
func Sources() []string {
	return this().Keys()
}

// CloseInstance 关闭数据库连接实例
func CloseInstance(source ...string) {
	if len(source) > 0 {
		for _, s := range source {
			_ = CloseDB(GetInstance(s))
		}
	}
}

// CloseDB 关闭数据库连接
func CloseDB(db *gorm.DB) error {
	if d, err := db.DB(); err != nil {
		return errorx.Wrap(err, "get DB failed")
	} else if err = d.Close(); err != nil {
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
