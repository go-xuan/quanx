package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/common/constx"
)

var _handler *Handler

func this() *Handler {
	if _handler == nil {
		panic("the redis handler has not been initialized, please check the relevant config")
	}
	return _handler
}

func AddClient(config *Config, client redis.UniversalClient) {
	if _handler == nil {
		_handler = &Handler{
			multi:   false,
			config:  config,
			client:  client,
			configs: make(map[string]*Config),
			clients: make(map[string]redis.UniversalClient),
		}
		return
	}
	_handler.multi = true
	_handler.configs[config.Source] = config
	_handler.clients[config.Source] = client
	if config.Source == constx.DefaultSource {
		_handler.config = config
		_handler.client = client
	}
}

// Handler redis连接句柄
type Handler struct {
	multi   bool
	config  *Config
	client  redis.UniversalClient
	configs map[string]*Config
	clients map[string]redis.UniversalClient
}

func (h *Handler) GetClient(source ...string) redis.UniversalClient {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clients[source[0]]; ok {
			return client
		}
	}
	return h.client
}

func (h *Handler) GetConfig(source ...string) *Config {
	if h.multi && len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
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
func GetClient(source ...string) redis.UniversalClient {
	return this().GetClient(source...)
}

// Ping 连接检查
func Ping(ctx context.Context, source ...string) (bool, error) {
	if result, err := this().GetClient(source...).Ping(ctx).Result(); err != nil || result != "PONG" {
		return false, err
	} else {
		return true, nil
	}
}
