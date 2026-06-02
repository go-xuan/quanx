package elasticx

import (
	"github.com/go-xuan/quanx/configx"
	"github.com/olivere/elastic/v7"
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
func GetInstance(source ...string) *elastic.Client {
	return GetClient(source...).GetClient()
}
