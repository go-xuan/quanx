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

// 服务启动器
type Engine struct {
	Config           *Config      // 应用配置
	GinEngine        *gin.Engine  // gin引擎
	GinRouterLoaders []RouterLoad // gin路由加载方法
	ConfigBinders    []EngineFunc // 配置加载方法
}

// 应用运行函数
type EngineFunc func()

// 加载gin路由
type RouterLoad func(router *gin.RouterGroup)

// 服务启动
func (e *Engine) Run(functions ...RouterLoad) {
	defer Recover()
	// 初始化配置
	e.InitConfig()
	// 执行配置绑定函数
	e.ExecConfigBinders()
	// 初始化日志
	logx.InitLogger(&e.Config.Log)
	// 初始化Nacos
	nacosx.InitNacosCTL(&e.Config.Nacos)
	// 加载Nacos配置
	loadNacosConfig(e.Config)
	// 注册Nacos服务
	registerNacosServer(e.Config.Server)
	// 初始化Gorm
	gormx.InitGormCTL(&e.Config.Database)
	// 初始化redis
	redisx.InitRedisCTL(&e.Config.Redis)
	// 启动
	e.AddRouterLoaders(functions...)
	e.StartGin()
}

// 初始化应用配置
func (e *Engine) InitConfig() {
	if e.Config == nil {
		var config Config
		if err := configor.New(&configor.Config{
			Debug:       true,
			Environment: configor.ENV(),
		}).Load(&config, "config.yaml"); err != nil {
			panic(err)
		}
		if ipx.GetWLANIP() != "" {
			config.Server.Host = ipx.GetWLANIP()
		}
		config.Log.Name = config.Server.Name
		e.Config = &config
	}
}

// 添加配置绑定函数
func (e *Engine) AddConfigBinders(funcs ...EngineFunc) {
	e.ConfigBinders = append(e.ConfigBinders, funcs...)
}

// 执行配置绑定函数
func (e *Engine) ExecConfigBinders() {
	if e.ConfigBinders != nil && len(e.ConfigBinders) > 0 {
		for _, binder := range e.ConfigBinders {
			binder()
		}
	} else {
		log.Warn("ExecConfigBinders is empty !")
	}
}

// 添加gin的路由加载函数
func (e *Engine) AddRouterLoaders(loads ...RouterLoad) {
	e.GinRouterLoaders = append(e.GinRouterLoaders, loads...)
}

// 执行gin的路由加载函数
func (e *Engine) ExecRouterLoaders(group *gin.RouterGroup) {
	if e.GinRouterLoaders != nil && len(e.GinRouterLoaders) > 0 {
		for _, loader := range e.GinRouterLoaders {
			loader(group)
		}
	} else {
		log.Warn("GinRouterLoaders is empty !")
	}
}

// 启动gin
func (e *Engine) StartGin() {
	if e.Config.Server.Env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	e.GinEngine = gin.New()
	e.GinEngine.Use(logx.LoggerToFile(), gin.Recovery())
	_ = e.GinEngine.SetTrustedProxies([]string{e.Config.Server.Host})
	// 注册根路由
	group := e.GinEngine.Group(e.Config.Server.Prefix)
	// 执行路由注册函数
	e.ExecRouterLoaders(group)
	log.Info("=== API接口请求地址: http://" + e.Config.Server.Host + ":" + e.Config.Server.Port)
	if err := e.GinEngine.Run(":" + e.Config.Server.Port); err != nil {
		panic(err)
	}
	log.Info("服务启动成功！！！")
}
