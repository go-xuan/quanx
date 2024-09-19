package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/app/constx"
)

var handler *Handler

// Handler redis控制器
type Handler struct {
	multi     bool // 是否多redis数据库
	config    *Redis
	configMap map[string]*Redis
	client    redis.UniversalClient
	clientMap map[string]redis.UniversalClient
}

func This() *Handler {
	if handler == nil {
		panic("the redis handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Client(source ...string) redis.UniversalClient {
	return This().GetClient(source...)
}

func Initialized() bool {
	return handler != nil
}

func Ping(client redis.UniversalClient) (bool, error) {
	if result, err := client.Ping(context.Background()).Result(); err != nil || result != "PONG" {
		return false, err
	} else {
		return true, nil
	}
}

func (h *Handler) GetClient(source ...string) redis.UniversalClient {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if client, ok := h.clientMap[source[0]]; ok {
			return client
		}
	}
	return h.client
}

func (h *Handler) GetConfig(source ...string) *Redis {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if conf, ok := h.configMap[source[0]]; ok {
			return conf
		}
	}
	return h.config
}
