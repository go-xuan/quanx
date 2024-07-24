package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/app/constx"
)

var handler *Handler

// Handler redis控制器
type Handler struct {
	Multi     bool // 是否多redis数据库
	Client    *redis.UniversalClient
	Config    *Redis
	clientMap map[string]*redis.UniversalClient
	ConfigMap map[string]*Redis
}

func This() *Handler {
	if !Initialized() {
		panic("the redis handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Initialized() bool {
	return handler != nil
}

func Ping(client redis.UniversalClient) (bool, error) {
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Client(source ...string) redis.UniversalClient {
	return This().GetClient(source...)
}

func (h *Handler) GetClient(source ...string) redis.UniversalClient {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if client, ok := h.clientMap[source[0]]; ok {
			return *client
		}
	}
	return *h.Client
}

func (h *Handler) GetConfig(source ...string) *Redis {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}
