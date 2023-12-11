package redisx

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var instance *Handler

// redis控制器
type Handler struct {
	Cmd       redis.Cmdable
	Config    *Config
	CmdMap    map[string]redis.Cmdable
	ConfigMap map[string]*Config
}

func This() *Handler {
	if instance == nil {
		panic("The redis instance has not been initialized, please check the relevant config")
	}
	return instance
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
		if cmd, ok := instance.CmdMap[source[0]]; ok {
			return cmd
		}
	}
	return instance.Cmd
}

func GetConfig(source ...string) *Config {
	if len(source) > 0 {
		if conf, ok := instance.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return instance.Config
}
