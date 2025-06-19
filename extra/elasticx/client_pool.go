package elasticx

import (
	"github.com/olivere/elastic/v7"

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
		panic("elastic client not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, cli *elastic.Client) {
	if config == nil || cli == nil {
		return
	}
	client := &Client{config, cli}
	if pool == nil {
		pool = enumx.NewStringEnum[*Client]()
		pool.Add(constx.DefaultSource, client)
	}
	pool.Add(config.Source, client)
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

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}

// GetInstance 获取数据库连接
func GetInstance(source ...string) *elastic.Client {
	return GetClient(source...).Instance()
}
