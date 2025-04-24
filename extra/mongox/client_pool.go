package mongox

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/types/enumx"
)

var pool *enumx.Enum[string, *Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

func this() *enumx.Enum[string, *Client] {
	if pool == nil {
		panic("mongo client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, client *mongo.Client) {
	if pool == nil {
		pool = enumx.NewStringEnum[*Client]()
	}
	pool.Add(config.Source, &Client{config, client})
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := pool.Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *mongo.Client {
	return GetClient(source...).Instance()
}

func GetDatabase(source ...string) *mongo.Database {
	if client := GetClient(source...); client != nil {
		return client.Instance().Database(client.Config().Database)
	}
	return nil
}
