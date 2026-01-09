package elasticx

import (
	"github.com/go-xuan/typex"
	"github.com/olivere/elastic/v7"

	"github.com/go-xuan/quanx/constx"
)

// 客户端池
var pool *typex.Enum[string, *Client]

// Pool 获取客户端池
func Pool() *typex.Enum[string, *Client] {
	if !Initialized() {
		panic("client pool not initialized")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

// AddClient 添加客户端
func AddClient(source string, client *Client) {
	if client == nil {
		return
	}
	if !Initialized() {
		pool = typex.NewStringEnum[*Client]()
		pool.Add(constx.DefaultSource, client)
	}
	pool.Add(source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) *Client {
	if len(source) > 0 && source[0] != "" {
		if client := Pool().Get(source[0]); client != nil {
			return client
		}
	}
	return Pool().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetESClient 获取数据库连接
func GetESClient(source ...string) *elastic.Client {
	return GetClient(source...).GetClient()
}
