package cachex

import (
	"github.com/go-xuan/quanx/common/constx"
)

var handler *Handler

type Handler struct {
	Multi     bool // 是否多缓存
	Client    *CacheClient
	ClientMap map[string]*CacheClient
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

func Client(name ...string) *CacheClient {
	return This().GetClient(name...)
}

func (h *Handler) GetClient(name ...string) *CacheClient {
	if len(name) > 0 && name[0] != constx.Default {
		if client, ok := h.ClientMap[name[0]]; ok {
			return client
		}
	}
	return h.Client
}
