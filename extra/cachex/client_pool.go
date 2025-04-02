package cachex

import (
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/types/enumx"
	"github.com/go-xuan/quanx/types/stringx"
)

var pool *enumx.Enum[string, Client]

func this() *enumx.Enum[string, Client] {
	if pool == nil {
		panic("cache cache not initialized, please check the relevant config")
	}
	return pool
}

func AddClient(config *Config, client Client) {
	if pool == nil {
		pool = enumx.NewStringEnum[Client]()
	}
	pool.Add(config.Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	return this().Get(stringx.Default(constx.DefaultSource, source...))
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}
