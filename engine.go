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
	"github.com/go-xuan/quanx/frame/confx"
	"github.com/go-xuan/quanx/frame/nacosx"
	"github.com/go-xuan/quanx/os/ipx"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/os/marshalx"
	"github.com/go-xuan/quanx/utils/anyx"
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
	config         *Config                        // 服务配置 使用 initConfig()将配置文件加载到此
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
	Server        *Server              `yaml:"server"`        // 服务配置
	Log           *logx.Log            `yaml:"log"`           // 日志配置
	Nacos         *nacosx.Nacos        `yaml:"nacos"`         // nacos访问配置
	Database      *gormx.Database      `yaml:"database"`      // 数据源配置
	Redis         *redisx.Redis        `yaml:"redis"`         // redis配置
	MultiDatabase *gormx.MultiDatabase `yaml:"multiDatabase"` // 多数据源配置
	MultiRedis    *redisx.MultiRedis   `yaml:"multiRedis"`    // 多redis配置
}

// gin路由加载器
type RouterLoader func(router *gin.RouterGroup)

// 自定义函数
type CustomFunc func()

// 服务启动标识
type Flag int

const (
	DisableGin              Flag = iota // 关闭gin
	EnableNacos                         // 启用nacos
	MultiDB                             // 启用多数据源
	MultiRedis                          // 启用多Redis
	InitConfigAlready                   // 已加载配置
	InitCommonAlready                   // 已初始化基础配置
	RunCustomFuncAlready                // 已执行自定义函数
	RunConfiguratorsAlready             // 已执行配置器
)

// 服务配置
type Server struct {
	Name   string `yaml:"name"`                     // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix" default:"app"`     // RESTFul api prefix（接口根路由）
	Debug  bool   `yaml:"debug"`                    // 是否调试环境
}

// 服务地址
func (s *Server) HttpUrl() string {
	return fmt.Sprintf(`http://%s:%d/%s`, s.Host, s.Port, strings.TrimPrefix(s.Prefix, "/"))
}

// 服务运行
func (e *Engine) RUN() {
	e.PREPARE()  // 服务准备
	e.STARTGIN() // 服务启动
}

// 服务准备
func (e *Engine) PREPARE() {
	// 加载配置
	if !e.flag[InitConfigAlready] {
		e.initConfig()
	}
	// 初始化基础配置(日志/Nacos/gorm/Redis)
	if !e.flag[InitCommonAlready] {
		e.initCommon()
	}
	// 执行配置器
	if !e.flag[RunConfiguratorsAlready] {
		e.runConfigurators()
	}
	// 执行自定义函数
	if !e.flag[RunCustomFuncAlready] {
		e.runCustomFunc()
	}
}

// 服务启动
func (e *Engine) STARTGIN() {
	if !e.flag[DisableGin] {
		defer PanicRecover()
		e.startGin()
	}
}

// 加载服务配置
func (e *Engine) initConfig() {
	var config = &Config{Server: &Server{}}
	// 读取本地配置
	if !e.flag[DisableGin] {
		var path = e.GetConfigPath("config.yaml")
		if err := marshalx.LoadFromFile(path, config); err != nil {
			log.Error("加载服务配置失败!")
			panic(err)
		}
		if config.Server.Host == "" {
			config.Server.Host = ipx.GetWLANIP()
		}
	}
	// 初始化nacos
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
	e.flag[InitConfigAlready] = true
}

// 初始化日志/Nacos/gorm/Redis
func (e *Engine) initCommon() {
	var serverName = anyx.IfZero(e.config.Server.Name, "app")
	// 初始化日志
	e.RunConfigurator(anyx.IfZero(e.config.Log, logx.New(serverName)), true)
	// 连接数据库
	if e.flag[MultiDB] {
		e.config.MultiDatabase = anyx.IfZero(e.config.MultiDatabase, &gormx.MultiDatabase{})
		e.RunConfigurator(e.config.MultiDatabase)
		if gormx.Initialized() && gormx.This().Multi {
			for source := range gormx.This().DBMap {
				if dst, ok := e.gormTables[source]; ok {
					if err := gormx.This().InitGormTable(source, dst...); err != nil {
						log.Error("failed to initialize the table structure !")
						panic(err)
					}
				}
			}
		}
	} else if !gormx.Initialized() {
		e.config.Database = anyx.IfZero(e.config.Database, &gormx.Database{})
		e.RunConfigurator(e.config.Database)
		if dst, ok := e.gormTables[constx.Default]; ok {
			if err := gormx.This().InitGormTable(constx.Default, dst...); err != nil {
				log.Error("failed to initialize the table structure !")
				panic(err)
			}
		}
	}
	// 连接redis
	if e.flag[MultiRedis] {
		e.config.MultiRedis = anyx.IfZero(e.config.MultiRedis, &redisx.MultiRedis{})
		e.RunConfigurator(e.config.MultiRedis)
	} else if !redisx.Initialized() {
		e.config.Redis = anyx.IfZero(e.config.Redis, &redisx.Redis{})
		e.RunConfigurator(e.config.Redis)
	}
	e.flag[InitCommonAlready] = true
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
	e.flag[RunCustomFuncAlready] = true
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
	e.flag[RunConfiguratorsAlready] = true
}

// 运行配置器
func (e *Engine) RunConfigurator(conf confx.Configurator[any], must ...bool) {
	var ok = anyx.Default(must, false)
	if reader := conf.Reader(); reader != nil {
		if e.flag[EnableNacos] {
			reader.NacosGroup = e.config.Server.Name
			if err := nacosx.NewConfig(reader.NacosGroup, reader.NacosDataId).LoadConfig(conf); err == nil {
				ok = true
			}
		} else {
			if err := marshalx.LoadFromFile(e.GetConfigPath(reader.FilePath), conf); err == nil {
				ok = true
			}
		}
	}
	if ok {
		if err := conf.Run(); err != nil {
			log.Error(conf.Title() + " error!")
			panic(err)
		}
		log.Info(conf.Title() + " completed!")
	}
}

// 初始化本地配置项（立即加载）
func (e *Engine) InitLocalConfig(config interface{}, filePath string) {
	if err := marshalx.LoadFromFile(filePath, config); err != nil {
		panic(err)
	}
}

// 初始化Nacos配置项（以自定义函数的形式延迟加载）
func (e *Engine) InitNacosConfig(config interface{}, dataId string, listen ...bool) {
	e.AddCustomFunc(func() {
		if nacosx.This().ConfigClient == nil {
			panic("nacos config client is uninitialized !")
		}
		var item = nacosx.NewConfig(e.config.Server.Name, dataId)
		if len(listen) > 0 {
			item.Listen = listen[0]
		}
		// 加载微服务配置
		if err := item.LoadConfig(config); err != nil {
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
func (e *Engine) GetConfigPath(name string) string {
	return anyx.IfZeroElse(e.configDir, name, filepath.Join(e.configDir, name))
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
