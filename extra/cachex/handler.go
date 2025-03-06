package cachex

import (
	"github.com/go-xuan/quanx/common/constx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the cache handler has not been initialized, please check the relevant config")
	}
	return _handler
}

type Handler struct {
	multi     bool // 是否多缓存
	client    Client
	clientMap map[string]Client
}

func (h *Handler) GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clientMap[source[0]]; ok {
			return client
		}
	}
	return h.client
}

// GetConfig 获取配置
func GetConfig(source ...string) *Config {
	return this().GetClient(source...).Config()
}

// GetClient 获取客户端
func GetClient(source ...string) Client {
	return this().GetClient(source...)
}
