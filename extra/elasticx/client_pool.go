package elasticx

import (
	"github.com/olivere/elastic/v7"

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
		panic("elastic client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, client *elastic.Client) {
	if pool == nil {
		pool = enumx.NewStringEnum[*Client]()
	}
	pool.Add(config.Source, &Client{config, client})
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
func GetInstance(source ...string) *elastic.Client {
	return GetClient(source...).Instance()
}
