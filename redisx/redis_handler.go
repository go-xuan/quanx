package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var handler *Handler

// redis控制器
type Handler struct {
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
	_, err := cmd.Ping(context.Background()).Result()
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func GetCmd(source ...string) redis.Cmdable {
	if len(source) > 0 {
		if cmd, ok := handler.CmdMap[source[0]]; ok {
			return cmd
		}
	}
	return handler.Cmd
}

func GetConfig(source ...string) *Redis {
	if len(source) > 0 {
		if conf, ok := handler.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return handler.Config
}
