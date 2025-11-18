package mongox

import (
	"github.com/go-xuan/typex"
	"go.mongodb.org/mongo-driver/mongo"
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
	return pool != nil
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
func AddClient(config *Config, cli *mongo.Client) {
	if config == nil || cli == nil {
		return
	}
	client := &Client{config, cli}
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		this().Add("default", client)
	}
	this().Add(config.Source, client)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *mongo.Client {
	return GetClient(source...).GetInstance()
}

func GetDatabase(source ...string) *mongo.Database {
	if client := GetClient(source...); client != nil {
		return client.GetInstance().Database(client.GetConfig().Database)
	}
	return nil
}
