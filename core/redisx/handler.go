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

func Initialized() bool {
	return _handler != nil
}

func GetConfig(source ...string) *Config {
	return this().GetConfig(source...)
}

func GetClient(source ...string) redis.UniversalClient {
	return this().GetClient(source...)
}

func Ping(ctx context.Context, source ...string) (bool, error) {
	if result, err := this().GetClient(source...).Ping(ctx).Result(); err != nil || result != "PONG" {
		return false, err
	} else {
		return true, nil
	}
}

// Handler redis连接句柄
type Handler struct {
	multi   bool // 是否多redis数据库
	config  *Config
	configs map[string]*Config
	client  redis.UniversalClient
	clients map[string]redis.UniversalClient
}

func (h *Handler) GetClient(source ...string) redis.UniversalClient {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if client, ok := h.clients[source[0]]; ok {
			return client
		}
	}
	return h.client
}

func (h *Handler) GetConfig(source ...string) *Config {
	if len(source) > 0 && source[0] != constx.DefaultSource {
		if conf, ok := h.configs[source[0]]; ok {
			return conf
		}
	}
	return h.config
}
