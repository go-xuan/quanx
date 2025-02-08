package mongox

import (
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-xuan/quanx/common/constx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the mongo handler has not been initialized, please check the relevant config")
	}
	return _handler
}

type Handler struct {
	multi   bool
	config  *Config
	client  *mongo.Client
	configs map[string]*Config
	clients map[string]*mongo.Client
}

func (h *Handler) GetConfig(source ...string) *Config {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
}

func (h *Handler) GetClient(source ...string) *mongo.Client {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clients[source[0]]; ok {
			return client
		}
	}
	return h.client
}

// IsInitialized 是否初始化
func IsInitialized() bool {
	return _handler != nil
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

// GetClient 获取客户端
func GetClient(source ...string) *mongo.Client {
	return this().GetClient(source...)
}

func GetDatabase(source ...string) *mongo.Database {
	if conf, client := GetConfig(source...), GetClient(source...); conf != nil && client != nil {
		return client.Database(conf.Database)
	}
	return nil
}
