package mongox

import (
	"github.com/go-xuan/typex"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-xuan/quanx/constx"
)

// 客户端池
var pool *typex.Enum[string, *Client]

// 获取客户端池
func this() *typex.Enum[string, *Client] {
	if !Initialized() {
		panic("mongo client pool not initialized, please check the relevant config")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.DefaultSource)
}

// AddClient 添加客户端
func AddClient(client *Client) {
	if client != nil {
		if !Initialized() {
			pool = typex.NewStringEnum[*Client]()
			pool.Add(constx.DefaultSource, client)
		}
		this().Add(client.GetConfig().Source, client)
	}
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *mongo.Client {
	return GetClient(source...).GetInstance()
}

// GetDatabase 获取数据库
func GetDatabase(source ...string) *mongo.Database {
	if client := GetClient(source...); client != nil {
		database := client.GetConfig().Database
		return client.GetInstance().Database(database)
	}
	return nil
}
