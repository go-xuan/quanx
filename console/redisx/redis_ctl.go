package redisx

import (
	"context"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var CTL *Controller

// redis控制器
type Controller struct {
	Cmd       redis.Cmdable
	Config    *Config
	CmdMap    map[string]redis.Cmdable
	ConfigMap map[string]*Config
}

// 初始化redis控制器
func Init(configs Configs) {
	if len(configs) == 0 {
		log.Error("redis配置为空!")
		return
	}
	CTL = &Controller{
		CmdMap:    make(map[string]redis.Cmdable),
		ConfigMap: make(map[string]*Config),
	}
	for i, conf := range configs {
		var cmd = conf.NewRedisCmdable()
		if ok, err := Ping(cmd); ok && err == nil {
			CTL.CmdMap[conf.Source] = cmd
			CTL.ConfigMap[conf.Source] = conf
			log.Info("redis连接成功! ", conf.Format())
			if i == 0 || conf.Source == "default" {
				CTL.Cmd = cmd
				CTL.Config = conf
			}
		} else {
			log.Error("redis连接失败! ", conf.Format())
			log.Error("error : ", err)
		}
	}
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
		if cmd, ok := CTL.CmdMap[source[0]]; ok {
			return cmd
		}
	}
	return CTL.Cmd
}

func GetConfig(source ...string) *Config {
	if len(source) > 0 {
		if conf, ok := CTL.ConfigMap[source[0]]; ok {
			return conf
		}
	}
	return CTL.Config
}
