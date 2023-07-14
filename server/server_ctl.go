package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/quanxiaoxuan/quanx/middleware/gormx"
	"github.com/quanxiaoxuan/quanx/middleware/logx"
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	"github.com/quanxiaoxuan/quanx/middleware/redisx"
	"github.com/quanxiaoxuan/quanx/utils/ipx"
	log "github.com/sirupsen/logrus"
)

// Gorm控制器
type Starter struct {
	Config      *AppConfig   // 应用配置
	GinEngine   *gin.Engine  // gin引擎
	RouterFuncs []RouterFunc // gin路由注册方法集

}

// 应用配置
type AppConfig struct {
	Server   Server                `yaml:"server"`   // 服务配置
	Log      logx.Config           `yaml:"log"`      // 日志配置
	Nacos    nacosx.Config         `yaml:"nacos"`    // nacos访问配置
	Configs  nacosx.ConfigItemList `yaml:"configs"`  // nacos配置清单
	Database gormx.Config          `yaml:"database"` // 数据库访问配置
	Redis    redisx.Config         `yaml:"redis"`    // redis配置
}

// 服务配置
type Server struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host" default:"127.0.0.1"`
	Port   string `yaml:"port" default:"8888"`
	Env    string `yaml:"env" default:"localhost"` // server run environment
	Prefix string `yaml:"prefix" default:"api"`    // prefix of api (RESTFul api prefix)
}

// gin框架路由注册方法
type RouterFunc func(router *gin.RouterGroup)

// 初始化应用配置
func (s *Starter) InitConfig() {
	if s.Config == nil {
		var config AppConfig
		if err := configor.New(&configor.Config{
			Debug:       true,
			Environment: configor.ENV(),
		}).Load(&config, DefaultConfig); err != nil {
			panic(err)
		}
		if ipx.GetWLANIP() != "" {
			config.Server.Host = ipx.GetWLANIP()
		}
		config.Log.Name = config.Server.Name
		s.Config = &config
	}
}

// 添加gin的路由注册方法
func (s *Starter) AddRouterFuncs(f ...RouterFunc) {
	s.RouterFuncs = append(s.RouterFuncs, f...)
}

// 添加gin的路由注册方法
func (s *Starter) RouterFuncsCheck() bool {
	if len(s.RouterFuncs) == 0 {
		return false
	}
	return true
}

// 启动gin
func (s *Starter) StartGin() {
	if s.Config.Server.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	s.GinEngine = gin.New()
	s.GinEngine.Use(logx.LoggerToFile(), gin.Recovery())
	_ = s.GinEngine.SetTrustedProxies([]string{s.Config.Server.Host})
	// 注册根路由
	group := s.GinEngine.Group(s.Config.Server.Prefix)
	// 执行路由注册函数

	if s.RouterFuncsCheck() {
		for _, f := range s.RouterFuncs {
			f(group)
		}
	} else {
		log.Warn("Starter.RouterFuncs is empty !")
	}
	log.Info("=== API接口请求地址: http://" + s.Config.Server.Host + ":" + s.Config.Server.Port)
	if err := s.GinEngine.Run(":" + s.Config.Server.Port); err != nil {
		panic(err)
	}
	log.Info("服务启动成功！！！")
}
