package cachex

import (
	"github.com/go-xuan/quanx/app/constx"
)

var handler *Handler

type Handler struct {
	Multi     bool // 是否多缓存
	Client    Client
	ClientMap map[string]Client
}

func This() *Handler {
	if !Initialized() {
		panic("the cache handler has not been initialized, please check the relevant config")
	}
	return handler
}

func Initialized() bool {
	return handler != nil
}

func GetClient(source ...string) Client {
	return This().GetClient(source...)
}

func (h *Handler) GetClient(source ...string) Client {
	if len(source) > 0 && source[0] != constx.DefaultKey {
		if client, ok := h.ClientMap[source[0]]; ok {
			return client
		}
	}
	return h.Client
}
