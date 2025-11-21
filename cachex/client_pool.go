package cachex

import (
	"github.com/go-xuan/typex"
	
	"github.com/go-xuan/quanx/constx"
)

var pool *typex.Enum[string, Client]

// 获取客户端池
func this() *typex.Enum[string, Client] {
	if !Initialized() {
		panic("cache client pool not initialized, please check the relevant config")
	}
	return pool
}

// Initialized 是否初始化
func Initialized() bool {
	return pool != nil
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.Default)
}

// AddClient 添加客户端
func AddClient(source string, client Client) {
	if !Initialized() {
		pool = typex.NewStringEnum[Client]()
		pool.Add(constx.Default, client)
	}
	this().Add(source, client)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).Config()
}
