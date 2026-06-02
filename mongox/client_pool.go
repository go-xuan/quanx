package mongox

import (
	"github.com/go-xuan/quanx/configx"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// 客户端池
var pool = configx.NewPool[*Client]()

// Pool 获取客户端池
func Pool() *configx.Pool[*Client] {
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool.Initialized()
}

// AddClient 添加客户端
func AddClient(source string, client *Client) {
	if client == nil {
		return
	}
	pool.Add(source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	return pool.Get(source...)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *mongo.Client {
	return GetClient(source...).GetClient()
}

// GetDatabase 获取数据库
func GetDatabase(source ...string) *mongo.Database {
	if client := GetClient(source...); client != nil {
		database := client.GetConfig().Database
		return client.GetClient().Database(database)
	}
	return nil
}

// Close 关闭所有数据库客户端
func Close() error {
	return pool.Close(func(client *Client) error {
		return client.Close()
	})
}
