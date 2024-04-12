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
	Cmd       redis.Cmdable
	Config    *Redis
	CmdMap    map[string]redis.Cmdable
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

func Ping(cmd redis.Cmdable) (bool, error) {
	if _, err := cmd.Ping(context.Background()).Result(); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func Cmdable(source ...string) redis.Cmdable {
	return This().GetCmdable(source...)
}

func (h *Handler) GetCmdable(source ...string) redis.Cmdable {
	if len(source) > 0 && source[0] != constx.Default {
		if cmd, ok := h.CmdMap[source[0]]; ok {
			return cmd
		}
	}
	return h.Cmd
}

func (h *Handler) GetConfig(source ...string) *Redis {
	if len(source) > 0 && source[0] != constx.Default {
		if conf, ok := h.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return h.Config
}
