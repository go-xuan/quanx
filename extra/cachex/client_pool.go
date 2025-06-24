package cachex

import (
	"github.com/go-xuan/typex"
)

var pool *typex.Enum[string, Client]

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil && pool.Len() > 0
}

func this() *typex.Enum[string, Client] {
	if pool == nil {
		panic("cache cache not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, client Client) {
	if pool == nil {
		pool = typex.NewStringEnum[Client]()
		pool.Add("default", client)
	}
	pool.Add(config.Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get("default")
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}
