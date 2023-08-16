package server

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
	"github.com/quanxiaoxuan/quanx/middleware/elasticx"
	"github.com/quanxiaoxuan/quanx/middleware/gormx"
	"github.com/quanxiaoxuan/quanx/middleware/hugegraphx"
	"github.com/quanxiaoxuan/quanx/middleware/logx"
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	"github.com/quanxiaoxuan/quanx/middleware/redisx"
	"github.com/quanxiaoxuan/quanx/utils/ipx"
	log "github.com/sirupsen/logrus"
)

var engine *Engine

// 初始化运行器
func NewEngine() *Engine {
	if engine == nil {
		engine = &Engine{}
	}
	return engine
}

// 服务启动器
type Engine struct {
	Config           *Config           // 应用配置
	BeforeFunc       EngineFunc        // 前置函数
	Middleware       []EngineFunc      // 中间件函数
	GinEngine        *gin.Engine       // gin引擎
	GinRouterLoaders []RouterLoader    // gin路由加载方法
	GinMiddleware    []gin.HandlerFunc // gin中间件
}

// 引擎函数
type EngineFunc func()

// gin路由加载函数
type RouterLoader func(router *gin.RouterGroup)

// 服务启动
func (e *Engine) RUN() {
	defer Recover()
	// 执行前置函数
	e.ExecBeforeFunc()
	// 初始化配置
	e.InitConfig()
	// 初始化中间件
	e.ExecMiddlewares()
	// 启动gin
	e.StartGin()
}

// 初始化配置
func (e *Engine) InitConfig() {
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

// 初始化日志
func (e *Engine) InitLogger() {
	logx.InitLogger(&e.Config.Log)
}

// 初始化Nacos
func (e *Engine) InitNacos() {
	if e.Config.Nacos.Address != "" {
		// 初始化Nacos
		nacosx.Init(&e.Config.Nacos)
		// 加载Nacos配置
		nacosx.LoadNacosConfig(e.Config.NacosOptions, e.Config.Server.Name, e.Config)
		// 注册Nacos服务
		nacosx.RegisterInstance(nacosx.ServerInstance{
			Group: e.Config.Server.Env,
			Name:  e.Config.Server.Name,
			Host:  e.Config.Server.Host,
			Port:  e.Config.Server.Port,
		})
	}
}

// 初始化Gorm
func (e *Engine) InitGorm() {
	if e.Config.Database.Host != "" {
		gormx.Init(&e.Config.Database)
	}
}

// 初始化Redis
func (e *Engine) InitRedis() {
	if e.Config.Redis.Host != "" {
		redisx.Init(&e.Config.Redis)
	}
}

// 初始化ElasticSearch
func (e *Engine) InitElasticSearch() {
	if e.Config.Elastic.Host != "" {
		elasticx.Init(&e.Config.Elastic)
	}
}

// 初始化hugegraph
func (e *Engine) InitHugegraph() {
	if e.Config.Hugegraph.Host != "" {
		hugegraphx.Init(&e.Config.Hugegraph)
	}
}

// 初始化hugegraph
func (e *Engine) InitMinio() {
	if e.Config.Minio.Host != "" {
		hugegraphx.Init(&e.Config.Hugegraph)
	}
}

// 执行前置函数
func (e *Engine) ExecBeforeFunc() {
	if e.BeforeFunc != nil {
		e.BeforeFunc()
	}
}

// 执行前置函数
func (e *Engine) ExecMiddlewares() {
	if e.Middleware != nil && len(e.Middleware) > 0 {
		for _, engineFunc := range e.Middleware {
			engineFunc()
		}
	}
}

// 添加中间件初始化函数
func (e *Engine) AddMiddleware(middleware ...EngineFunc) {
	e.Middleware = append(e.Middleware, middleware...)
}

// 执行gin的路由加载函数
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

// 获取配置
func (e *Engine) GetConfig(config *Config) {
	config = e.Config
}

// 服务保活
func Recover() {
	if err := recover(); err != nil {
		log.Errorf("服务启动失败！错误为 : %s", err)
		return
	}
	select {}
}
