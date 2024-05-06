package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/go-xuan/quanx/common/constx"
)

var handler *Handler

// redis控制器
type Handler struct {
	Multi     bool // 是否多redis数据库
	Client    *redis.UniversalClient
	Config    *Redis
	clientMap map[string]*redis.UniversalClient
	ConfigMap map[string]*Redis
}

func This() *Handler {
	if !Initialized() {
		panic("The redis handler has not been initialized, please check the relevant config")
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

func Client(name ...string) redis.UniversalClient {
	return This().GetClient(name...)
}

func (h *Handler) GetClient(name ...string) redis.UniversalClient {
	if len(name) > 0 && name[0] != constx.Default {
		if client, ok := h.clientMap[name[0]]; ok {
			return *client
		}
	}
	return *h.Client
}

func (h *Handler) GetConfig(name ...string) *Redis {
	if len(name) > 0 && name[0] != constx.Default {
		if conf, ok := h.ConfigMap[name[0]]; ok {
			return conf
		}
	}
	return h.Config
}
