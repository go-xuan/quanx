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

func AddClient(config *Config, client Client) {
	if _handler == nil {
		_handler = &Handler{
			multi:   false,
			client:  client,
			clients: make(map[string]Client),
		}
		return
	}
	_handler.multi = true
	_handler.clients[config.Source] = client
	if config.Source == constx.DefaultSource {
		_handler.config = config
		_handler.client = client
	}
}

type Handler struct {
	multi   bool // 是否多缓存
	client  Client
	config  *Config
	configs map[string]*Config
	clients map[string]Client
}

func (h *Handler) GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clients[source[0]]; ok {
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
