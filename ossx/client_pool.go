package ossx

import (
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/typex"
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
	return pool != nil && pool.Len() > 0
}

// AddClient 添加客户端
func AddClient(client Client) {
	if !Initialized() {
		pool = typex.NewStringEnum[Client]()
		pool.Add(constx.DefaultSource, client)
	}
	this().Add(client.GetConfig().Source, client)
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != "" {
		if client := this().Get(source[0]); client != nil {
			return client
		}
	}
	return this().Get(constx.DefaultSource)
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return GetClient(source...).GetConfig()
}

// GetInstance 获取实例
func GetInstance[INS any](source ...string) INS {
	instance := GetClient(source...).GetInstance()
	if ins, ok := instance.(INS); ok {
		return ins
	}
	var ins INS
	return ins
}
