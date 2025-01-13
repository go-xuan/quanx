package cachex

import (
	"github.com/go-xuan/quanx/common/constx"
)

var _handler *Handler

type Handler struct {
	multi     bool // 是否多缓存
	client    Client
	clientMap map[string]Client
}

func this() *Handler {
	if _handler == nil {
		panic("the cache handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func GetConfig(source ...string) *Config {
	return this().GetClient(source...).Config()
}

func GetClient(source ...string) Client {
	return this().GetClient(source...)
}

func (h *Handler) GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clientMap[source[0]]; ok {
			return client
		}
	}
	return h.client
}
