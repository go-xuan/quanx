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

func GetClient(source ...string) *CacheClient {
	return This().GetClient(source...)
}

func (h *Handler) GetClient(source ...string) *CacheClient {
	if len(source) > 0 && source[0] != constx.DefaultSourceName {
		if client, ok := h.ClientMap[source[0]]; ok {
			return client
		}
	}
	return h.Client
}
