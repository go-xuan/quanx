package engine

import (
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/public/gormx"
	"github.com/go-xuan/quanx/public/logx"
	"github.com/go-xuan/quanx/public/nacosx"
	"github.com/go-xuan/quanx/public/redisx"
	"github.com/go-xuan/quanx/utils/ipx"
	"github.com/go-xuan/quanx/utils/structx"
)

var engine *Engine

// 初始化启动器
func GetEngine() *Engine {
	if engine == nil {
		engine = &Engine{
			ConfigPath: "config.yaml",
			GormModel:  make(map[string][]interface{}),
		}
	}
	return engine
}

// 服务启动器
type Engine struct {
	Config           *Config                  // 服务依赖配置
	ConfigPath       string                   // 配置文件路径
	Initializers     []Initializer            // 服务初始化方法
	GinEngine        *gin.Engine              // gin引擎
	GinRouterLoaders []RouterLoader           // gin路由加载方法
	GinMiddleware    []gin.HandlerFunc        // gin中间件
	GormModel        map[string][]interface{} // gorm初始化model
}

// 初始化方法
type Initializer func()

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 服务启动
func (e *Engine) RUN() {
	defer Recover()
	// 加载配置
	e.LoadConfig()
	// 初始化（日志、nacos、gorm、redis）
	e.InitServer()
	// 执行初始化方法
	e.ExecInitializers()
	// 启动gin
	e.StartGin()
}

func (e *Engine) SetConfigPath(path string) {
	e.ConfigPath = path
}

// 加载配置
func (e *Engine) LoadConfig() {
	if e.ConfigPath == "" {
		e.ConfigPath = "config.yaml"
	}
	var config = &Config{}
	if err := structx.ReadConfigToPointer(e.ConfigPath, config); err != nil {
		log.Error("加载应用配置失败！")
		panic(err)
	}
	if ipx.GetWLANIP() != "" {
		config.Server.Host = ipx.GetWLANIP()
	}
	e.Config = config
}

// 初始化服务基本
func (e *Engine) InitServer() {
	// 初始化日志
	logx.InitLogger(e.Config.Log, e.Config.Server.Name)
	// 初始化Nacos
	if e.Config.Nacos != nil {
		// 初始化Nacos
		nacosx.Init(e.Config.Nacos)
		// 加载Nacos配置
		nacosx.LoadNacosConfig(e.Config.Nacos.Options, e.Config)
		// 注册Nacos服务
		nacosx.RegisterInstance(
			nacosx.ServerInstance{
				Group: e.Config.Server.Env,
				Name:  e.Config.Server.Name,
				Host:  e.Config.Server.Host,
				Port:  e.Config.Server.Port,
			})
	}
	// 初始化Gorm
	if len(e.Config.Database) > 0 {
		gormx.Init(e.Config.Database)
		for _, item := range e.Config.Database {
			if models, ok := e.GormModel[item.Source]; ok {
				if err := gormx.CTL.InitTable(item.Source, models...); err != nil {
					panic(err)
				}
			}
		}
	}
	// 初始化Redis
	if len(e.Config.Redis) > 0 {
		redisx.Init(e.Config.Redis)
	}
}

// 添加初始化方法
func (e *Engine) AddInitializer(init ...Initializer) {
	e.Initializers = append(e.Initializers, init...)
}

// 执行初始化方法
func (e *Engine) ExecInitializers() {
	if e.Initializers != nil && len(e.Initializers) > 0 {
		for _, engineFunc := range e.Initializers {
			engineFunc()
		}
	}
}

// 添加gorm初始化的model模型
func (e *Engine) AddModel(dst ...interface{}) {
	if len(dst) > 0 {
		var source = "default"
		if reflect.TypeOf(dst[0]).Kind() == reflect.String {
			source = dst[0].(string)
			dst = append(dst[:1], dst[2:]...)
		}
		e.AddSourceModel(source, dst...)
	}
}

// 添加gorm初始化的model模型
func (e *Engine) AddSourceModel(source string, dst ...interface{}) {
	e.GormModel[source] = append(e.GormModel[source], dst...)
}

// 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	e.GinMiddleware = append(e.GinMiddleware, middleware...)
}

// 添加gin的路由加载函数
func (e *Engine) AddRouterLoaders(load ...RouterLoader) {
	e.GinRouterLoaders = append(e.GinRouterLoaders, load...)
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
	if e.GinMiddleware != nil && len(e.GinMiddleware) > 0 {
		e.GinEngine.Use(e.GinMiddleware...)
	}
	_ = e.GinEngine.SetTrustedProxies([]string{e.Config.Server.Host})
	// 注册根路由，并执行路由注册函数
	var group = e.GinEngine.Group(e.Config.Server.Prefix)
	e.ExecRouterLoaders(group)
	var port = strconv.Itoa(e.Config.Server.Port)
	if err := e.GinEngine.Run(":" + port); err != nil {
		log.Error("gin-Engine 运行失败！！！")
		panic(err)
	}
	log.Info("服务启动成功！！！")
	log.Infof("API接口请求地址: http://%s:%s ", e.Config.Server.Host, port)
}

// 加载Nacos配置项
func (e *Engine) LoadNacosConfig(config interface{}, group string, dataId ...string) {
	if e.Config != nil && len(e.Config.Nacos.Options) > 0 {
		var optMap = make(map[string]nacosx.Option)
		for _, option := range e.Config.Nacos.Options {
			if option.Group == group {
				optMap[option.DataId] = option
			}
		}
		for _, item := range dataId {
			_ = optMap[item].LoadNacosConfig(config)
		}
	}
}

// 服务保活
func Recover() {
	if err := recover(); err != nil {
		log.Errorf("服务启动失败！错误为 : %s", err)
		return
	}
	select {}
}
