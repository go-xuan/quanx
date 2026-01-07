package dbx

import (
	"github.com/go-xuan/typex"
	"github.com/go-xuan/utilx/errorx"
	"gorm.io/gorm"

	"github.com/go-xuan/quanx/constx"
)

// 客户端池
var pool *typex.Enum[string, Client]

// 获取客户端池
func this() *typex.Enum[string, Client] {
	if !Initialized() {
		panic("gorm client pool not initialized, please check the relevant config")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

// Sources 所有数据源
func Sources() []string {
	return this().Keys()
}

// AddClient 添加客户端
func AddClient(client Client) {
	if !Initialized() {
		pool = typex.NewStringEnum[Client]()
		pool.Add(constx.DefaultSource, client)
	}
	this().Add(client.GetConfig().Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取实例
func GetInstance[INS any](source ...string) INS {
	instance := GetClient(source...).GetInstance()
	if ins, ok := instance.(INS); ok {
		return ins
	}
	var ins INS
	return ins
}

// GetGormDB 获取数据库连接
func GetGormDB(source ...string) *gorm.DB {
	return GetInstance[*gorm.DB](source...)
}

// CopyClient 复制数据库客户端
func CopyClient(source, target, database string) error {
	if source != "" && target != "" {
		sourceClient := GetClient(source)
		if sourceClient == nil {
			return errorx.New("source database client not exist")
		}
		targetClient, err := sourceClient.Copy(target, database)
		if err != nil {
			return errorx.Wrap(err, "copy database client failed")
		}
		this().Add(target, targetClient)
	}
	return nil
}

// Close 关闭所有数据库客户端
func Close() error {
	var err error
	this().Range(func(source string, client Client) bool {
		if err = client.Close(); err != nil {
			return true
		}
		return false
	})
	return err
}
