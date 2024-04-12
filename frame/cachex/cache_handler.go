package cachex

import (
	"context"
	"github.com/go-xuan/quanx/common/constx"
	"time"
)

var handler *Handler

type Handler struct {
	Multi     bool // 是否多缓存
	Client    *CacheClient[any]
	Config    *Cache
	ClientMap map[string]*CacheClient[any]
	ConfigMap map[string]*Cache
}

type CacheClient[T any] struct {
	Set   func(context.Context, string, T, time.Duration) // 设置缓存
	Get   func(context.Context, string) T                 // 获取缓存
	Del   func(context.Context, string)                   // 删除缓存
	Exist func(context.Context, string) bool              // 是否存在
}

func This() *Handler {
	if !Initialized() {
		panic("The cache handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Initialized() bool {
	return handler != nil
}

func Client(source ...string) *CacheClient[any] {
	return This().GetClient(source...)
}

func (h *Handler) GetClient(source ...string) *CacheClient[any] {
	if len(source) > 0 && source[0] != constx.Default {
		if client, ok := h.ClientMap[source[0]]; ok {
			return client
		}
	}
	return h.Client
}

func (h *Handler) GetConfig(source ...string) *Cache {
	if len(source) > 0 && source[0] != constx.Default {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}
