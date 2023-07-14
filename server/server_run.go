package server

import (
	log "github.com/sirupsen/logrus"

	"github.com/quanxiaoxuan/quanx/middleware/gormx"
	"github.com/quanxiaoxuan/quanx/middleware/logx"
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	"github.com/quanxiaoxuan/quanx/middleware/redisx"
)

const DefaultConfig = "config.yaml"

func New() *Starter {
	return &Starter{}
}

// 服务启动
func (s *Starter) Start(functions ...RouterFunc) {
	defer Recover()
	// 初始化配置
	s.InitConfig()
	// 初始化日志
	logx.InitLogger(&s.Config.Log)
	// 初始化Nacos
	nacosx.InitNacosCTL(&s.Config.Nacos)
	// 加载Nacos配置
	loadNacosConfig(s.Config)
	// 注册Nacos服务
	registerNacosServer(s.Config.Server)
	// 初始化Gorm
	gormx.InitGormCTL(&s.Config.Database)
	// 初始化redis
	redisx.InitRedisCTL(&s.Config.Redis)
	// 启动
	s.AddRouterFuncs(functions...)
	s.StartGin()
}

// 加载Nacos配置
func loadNacosConfig(appConf *AppConfig) {
	if nacosx.CTL.ConfigClient != nil {
		nacosx.LoadNacosConfig(appConf.Configs, appConf.Server.Name, &appConf)
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
