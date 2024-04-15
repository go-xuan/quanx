package quanx

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/db/gormx"
	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/frame/cachex"
	"github.com/go-xuan/quanx/frame/confx"
	"github.com/go-xuan/quanx/frame/nacosx"
	"github.com/go-xuan/quanx/os/ipx"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/go-xuan/quanx/utils/anyx"
	"github.com/go-xuan/quanx/utils/stringx"
)

var engine *Engine

// 初始化Engine
func GetEngine(modes ...Flag) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			configDir:      constx.ConfDir,
			customFuncs:    make([]CustomFunc, 0),
			configurators:  make([]confx.Configurator[any], 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			gormTables:     make(map[string][]gormx.Tabler[any]),
			flag:           make(map[Flag]bool),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.SetMode(modes...)
	}
	return engine
}

// 获取服务配置
func GetServer() *Server {
	return GetEngine().config.Server
}

// 服务配置器
type Engine struct {
	flag           map[Flag]bool                  // 服务运行标识
	config         *Config                        // 服务配置 使用 loadingConfig()将配置文件加载到此
	configDir      string                         // 服务配置文件文件夹, 使用 SetConfigDir()设置配置文件路径
	configurators  []confx.Configurator[any]      // 配置器，使用 AddConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	customFuncs    []CustomFunc                   // 自定义初始化函数 使用 AddCustomFunc()添加自定义函数
	ginEngine      *gin.Engine                    // gin框架实例
	ginLoaders     []RouterLoader                 // gin路由的预加载方法，使用 AddGinRouter()添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc              // gin中间件的预加载方法，使用 AddGinRouter()添加gin中间件
	gormTables     map[string][]gormx.Tabler[any] // gorm表结构对象，使用 AddTable()/ AddSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
}

// 服务配置
type Config struct {
	Server   *Server              `yaml:"server"`   // 服务配置
	Log      *logx.LogConfig      `yaml:"log"`      // 日志配置
	Nacos    *nacosx.NacosConfig  `yaml:"nacos"`    // nacos访问配置
	Database *gormx.MultiDatabase `yaml:"database"` // 数据源配置
	Redis    *redisx.MultiRedis   `yaml:"redis"`    // redis配置
	Cache    *cachex.MultiCache   `yaml:"cache"`    // 缓存配置
}

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 自定义函数
type CustomFunc func()

// 服务启动标识
type Flag int

const (
	Lightweight         Flag = iota // 轻量化
	EnableNacos                     // 启用nacos
	HasLoadingConfig                // 已加载配置
	HasBuildFrameBasic              // 已构建框架基础
	HasRunCustomFunc                // 已执行自定义函数
	HasRunConfigurators             // 已执行配置器
	MultiDatabase                   // 开启多数据源
	MultiRedis                      // 开启多redis源
	MultiCache                      // 开启多缓存源
)

// 服务配置
type Server struct {
	Name   string `yaml:"name"`                     // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix" default:"app"`     // api prefix（接口根路由）
	Debug  bool   `yaml:"debug"`                    // 是否调试环境
}

// 服务地址
func (s *Server) HttpUrl() string {
	return fmt.Sprintf(`http://%s:%d/%s`, s.Host, s.Port, strings.TrimPrefix(s.Prefix, "/"))
}

// 服务运行
func (e *Engine) RUN() {
	e.Prepare() // 服务准备
	e.Start()   // 服务启动
}

// 服务准备
func (e *Engine) Prepare() {
	// 加载配置
	if !e.flag[HasLoadingConfig] {
		e.loadingConfig()
	}
	// 构建框架基础（log/nacos/gorm/redis/cache）
	if !e.flag[HasBuildFrameBasic] {
		e.buildFrameBasic()
	}
	// 执行配置器
	if !e.flag[HasRunConfigurators] {
		e.runConfigurators()
	}
	// 执行自定义函数
	if !e.flag[HasRunCustomFunc] {
		e.runCustomFunc()
	}
}

// 服务启动
func (e *Engine) Start() {
	if !e.flag[Lightweight] {
		defer PanicRecover()
		e.startGin()
	}
}

// 加载服务配置
func (e *Engine) loadingConfig() {
	var config = &Config{Server: &Server{}}
	if !e.flag[Lightweight] {
		var path = e.GetConfigPath(constx.Config)
		if err := marshalx.UnmarshalFromFile(path, config); err != nil {
			log.Error("loading config.yaml failed!")
			panic(err)
		}
		if config.Server.Host == "" {
			config.Server.Host = ipx.GetWLANIP()
		}
	}
	// 从nacos加载配置
	if e.flag[EnableNacos] && config.Nacos != nil {
		e.RunConfigurator(config.Nacos, true)
		if config.Nacos.EnableNaming() {
			// 注册nacos服务Nacos
			nacosx.RegisterInstance(
				nacosx.ServerInstance{
					Name:  config.Server.Name,
					Host:  config.Server.Host,
					Port:  config.Server.Port,
					Group: config.Nacos.NameSpace,
				},
			)
		}
	}
	e.config = config
	e.flag[HasLoadingConfig] = true
}

// 构建框架基础（log/nacos/gorm/redis）
func (e *Engine) buildFrameBasic() {
	// 初始化日志
	var serverName = stringx.IfZero(e.config.Server.Name, "app")
	e.RunConfigurator(anyx.IfZero(e.config.Log, logx.New(serverName)), true)

	// 初始化数据库连接
	if e.flag[MultiDatabase] {
		e.config.Database = anyx.IfZero(e.config.Database, &gormx.MultiDatabase{})
		e.RunConfigurator(e.config.Database)
	} else {
		var database = &gormx.Database{}
		e.RunConfigurator(database)
		e.config.Database = &gormx.MultiDatabase{database}
	}
	// 初始化表结构
	if gormx.Initialized() {
		for source := range gormx.This().DBMap {
			if dst, ok := e.gormTables[source]; ok {
				if err := gormx.This().InitGormTable(source, dst...); err != nil {
					log.Error("create table failed!")
					panic(err)
				}
			}
		}
	}

	// 初始化redis连接
	if e.flag[MultiRedis] {
		e.config.Redis = anyx.IfZero(e.config.Redis, &redisx.MultiRedis{})
		e.RunConfigurator(e.config.Redis)
	} else {
		var redis = &redisx.Redis{}
		e.RunConfigurator(redis)
		e.config.Redis = &redisx.MultiRedis{redis}
	}

	// 初始化缓存
	if redisx.Initialized() {
		e.config.Cache = anyx.IfZero(e.config.Cache, &cachex.MultiCache{})
		e.RunConfigurator(e.config.Cache, true)

		if e.flag[MultiCache] {
			e.config.Cache = anyx.IfZero(e.config.Cache, &cachex.MultiCache{})
			e.RunConfigurator(e.config.Cache)
		} else {
			var cache = &cachex.Cache{}
			e.RunConfigurator(cache, true)
			e.config.Cache = &cachex.MultiCache{cache}
		}
	}
	// 完成框架基础初始化
	e.flag[HasBuildFrameBasic] = true
}

// 添加自定义函数
func (e *Engine) AddCustomFunc(starter ...CustomFunc) {
	if len(starter) > 0 {
		e.customFuncs = append(e.customFuncs, starter...)
	}
}

// 执行自定义函数
func (e *Engine) runCustomFunc() {
	if e.customFuncs != nil && len(e.customFuncs) > 0 {
		for _, f := range e.customFuncs {
			f()
		}
	}
	// 完成自定义函数执行
	e.flag[HasRunCustomFunc] = true
}

// 添加配置器
func (e *Engine) AddConfigurator(conf ...confx.Configurator[any]) {
	if len(conf) > 0 {
		e.configurators = append(e.configurators, conf...)
	}
}

// 加载配置器
func (e *Engine) runConfigurators() {
	if e.configurators != nil && len(e.configurators) > 0 {
		for _, config := range e.configurators {
			e.RunConfigurator(config)
		}
	}
	e.flag[HasRunConfigurators] = true
}

// 运行配置器
func (e *Engine) RunConfigurator(conf confx.Configurator[any], must ...bool) {
	var ok = anyx.Default(must, false)
	if reader := conf.Reader(); reader != nil {
		if e.flag[EnableNacos] {
			reader.NacosGroup = e.config.Server.Name
			if err := nacosx.NewConfig(reader.NacosGroup, reader.NacosDataId).Loading(conf); err == nil {
				ok = true
			}
		} else {
			if err := marshalx.UnmarshalFromFile(e.GetConfigPath(reader.FilePath), conf); err == nil {
				ok = true
			}
		}
	}
	if ok {
		if err := conf.Run(); err != nil {
			log.Error(conf.Theme(), " initialized failed!")
			panic(err)
		}
		log.Info(conf.Theme(), " initialized completed!")
	}
}

// 初始化本地配置项（立即加载）
func (e *Engine) LoadingLocalConfig(v any, path string) {
	if err := marshalx.UnmarshalFromFile(path, v); err != nil {
		panic(err)
	}
}

// 初始化Nacos配置项（以自定义函数的形式延迟加载）
func (e *Engine) LoadingNacosConfig(v any, dataId string, listen ...bool) {
	e.AddCustomFunc(func() {
		if nacosx.This().ConfigClient == nil {
			panic("nacos config client is uninitialized !")
		}
		var config = nacosx.NewConfig(e.config.Server.Name, dataId)
		config.Listen = anyx.Default(listen, false)
		// 加载微服务配置
		if err := config.Loading(v); err != nil {
			panic("loading nacos config failed : " + err.Error())
		}
	})
}

// 启动gin
func (e *Engine) startGin() {
	if e.config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	}
	if e.ginEngine == nil {
		e.ginEngine = gin.New()
	}
	e.ginEngine.Use(gin.Recovery(), logx.LoggerToFile())
	e.ginEngine.Use(e.ginMiddlewares...)
	_ = e.ginEngine.SetTrustedProxies([]string{e.config.Server.Host})
	// 注册服务根路由，并执行路由注册函数
	var group = e.ginEngine.Group(e.config.Server.Prefix)
	e.InitGinLoader(group)
	var port = ":" + strconv.Itoa(e.config.Server.Port)
	log.Info("API接口请求地址: http://" + e.config.Server.Host + port)
	if err := e.ginEngine.Run(port); err != nil {
		log.Error("gin-Engine run failed !")
		panic(err)
	}
}

// 设置模式
func (e *Engine) SetMode(flags ...Flag) {
	if len(flags) > 0 {
		for _, flag := range flags {
			e.flag[flag] = true
		}
	}
}

// 设置配置文件
func (e *Engine) SetConfigDir(dir string) {
	e.configDir = dir
}

// 设置配置文件
func (e *Engine) GetConfigPath(path string) string {
	if e.configDir == "" {
		return path
	} else {
		return filepath.Join(e.configDir, path)
	}
}

// 添加需要初始化的 gormx.Tabler 模型
func (e *Engine) AddTable(dst ...gormx.Tabler[any]) {
	e.AddSourceTable(constx.Default, dst...)
}

// 添加需要某个数据源的gormx.Table模型
func (e *Engine) AddSourceTable(source string, dst ...gormx.Tabler[any]) {
	if len(dst) > 0 {
		e.gormTables[source] = append(e.gormTables[source], dst...)
	}
}

// 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	if len(middleware) > 0 {
		e.ginMiddlewares = append(e.ginMiddlewares, middleware...)
	}
}

// 添加gin的路由加载函数
func (e *Engine) AddGinRouter(load ...RouterLoader) {
	if len(load) > 0 {
		e.ginLoaders = append(e.ginLoaders, load...)
	}
}

// 执行gin的路由加载函数
func (e *Engine) InitGinLoader(group *gin.RouterGroup) {
	if e.ginLoaders != nil && len(e.ginLoaders) > 0 {
		for _, loader := range e.ginLoaders {
			loader(group)
		}
	} else {
		log.Warn("engine.ginLoaders is empty !")
	}
}

// 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Error("server run panic : ", err)
		return
	}
	select {}
}
