package server

import (
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	log "github.com/sirupsen/logrus"
)

var Runner *Engine

// 初始化运行器
func NewEngine() *Engine {
	if Runner != nil {
		return Runner
	}
	return &Engine{}
}

// 获取服务配置
func GetConfig() *Config {
	if Runner != nil && Runner.Config != nil {
		return Runner.Config
	}
	return nil
}

// 加载Nacos配置
func loadNacosConfig(conf *Config) {
	if nacosx.CTL.ConfigClient != nil {
		nacosx.LoadNacosConfig(conf.Configs, conf.Server.Name, conf)
	} else {
		log.Error("未初始化nacos配置中心客户端!")
	}
}

// 注册Nacos服务
func registerNacosServer(server Server) {
	if nacosx.CTL.NamingClient != nil {
		nacosx.RegisterInstance(nacosx.ServerConfig{
			Group: server.Env,
			Name:  server.Name,
			Ip:    server.Host,
			Port:  server.Port,
		})
	} else {
		log.Error("未初始化nacos服务发现客户端!")
	}
}

// 服务保活
func Recover() {
	if err := recover(); err != nil {
		log.Errorf("服务启动失败，错误为 : %s", err)
		return
	}
	select {}
}
