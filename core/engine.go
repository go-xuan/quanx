package core

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/confx"
	"github.com/go-xuan/quanx/core/nacosx"
	"github.com/go-xuan/quanx/db/gormx"
	"github.com/go-xuan/quanx/db/redisx"
	"github.com/go-xuan/quanx/net/ipx"
	"github.com/go-xuan/quanx/os/cachex"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/os/syncx"
	"github.com/go-xuan/quanx/os/taskx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// 服务启动标识
type Flag int

const (
	NonGin        Flag = iota // 非gin项目
	EnableNacos               // 启用nacos
	MultiDatabase             // 开启多数据源
	MultiRedis                // 开启多redis源
	MultiCache                // 开启多缓存源
	UseQueue                  // 使用队列任务启动
)

// Init Queue Function Name
const (
	LoadingConfig    = "loading_config"    // 加载配置文件
	InitServerBasic  = "init_server_basic" // 初始化服务基础组件（log/nacos/gorm/redis/cache）
	RunConfigurators = "run_configurators" // 运行自定义配置器
	RunCustomFuncs   = "run_custom_funcs"  // 运行自定义函数
	StartServer      = "start_server"      // 服务启动
)

var engine *Engine

// 服务配置器
type Engine struct {
	flag           map[Flag]bool                  // 服务运行标识
	config         *Config                        // 服务配置数据，使用 loadingConfig()将配置文件加载到此
	configDir      string                         // 服务配置文件夹, 使用 SetConfigDir()设置配置文件读取路径
	configurators  []confx.Configurator           // 配置器，使用 AddConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	customFuncs    []func()                       // 自定义初始化函数 使用 AddCustomFunc()添加自定义函数
	ginEngine      *gin.Engine                    // gin框架引擎实例
	ginLoaders     []RouterLoader                 // gin路由的预加载方法，使用 AddGinRouter()添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc              // gin中间件的预加载方法，使用 AddGinRouter()添加gin中间件
	gormTables     map[string][]gormx.Tabler[any] // gorm表结构对象，使用 AddTable() / AddSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
	initQueue      *taskx.QueueScheduler          // Engine服务启动时的队列任务
}

// gin路由加载器
type RouterLoader func(*gin.RouterGroup)

// 获取服务配置
func GetServer() *Server {
	return GetEngine().config.Server
}

// 初始化Engine
func GetEngine(modes ...Flag) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			configDir:      constx.DefaultConfDir,
			customFuncs:    make([]func(), 0),
			configurators:  make([]confx.Configurator, 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			gormTables:     make(map[string][]gormx.Tabler[any]),
			flag:           make(map[Flag]bool),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.SetMode(modes...)
	}
	log.SetOutput(os.Stdout)
	log.SetFormatter(logx.DefaultFormatter())
	if engine.flag[UseQueue] {
		queue := taskx.Queue()
		queue.Add(LoadingConfig, engine.loadingConfig)       // 加载配置文件
		queue.Add(InitServerBasic, engine.InitServerBasic)   // 初始化服务基础组件（log/nacos/gorm/redis/cache）
		queue.Add(RunConfigurators, engine.runConfigurators) // 运行自定义配置器
		queue.Add(RunCustomFuncs, engine.runCustomFunc)      // 运行自定义函数
		queue.Add(StartServer, engine.startServer)           // 服务启动
		engine.initQueue = queue
	}
	return engine
}

// 服务运行
func (e *Engine) RUN() {
	if engine.flag[UseQueue] {
		// 任务队列方式启动
		engine.initQueue.Execute()
	} else {
		// 默认启动
		syncx.OnceDo(e.loadingConfig)    // 加载配置文件
		syncx.OnceDo(e.InitServerBasic)  // 初始化服务基础组件（log/nacos/gorm/redis/cache）
		syncx.OnceDo(e.runConfigurators) // 运行自定义配置器
		syncx.OnceDo(e.runCustomFunc)    // 运行自定义函数
		syncx.OnceDo(e.startServer)      // 服务启动
	}
}

// 加载服务配置
func (e *Engine) loadingConfig() {
	var config = &Config{Server: &Server{}}
	if !e.flag[NonGin] {
		var path = e.GetConfigPath(constx.DefaultServerConfig)
		if err := marshalx.UnmarshalFromFile(path, config); err != nil {
			log.Errorf("Loading %s Failed", path)
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
			if err := config.Nacos.Register(config.Server.Instance()); err != nil {
				panic(err)
			}
		}
	}
	e.config = config
}

// 构建框架基础（log/nacos/gorm/redis）
func (e *Engine) InitServerBasic() {
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
					log.Error("Create Table Failed")
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
		if e.flag[MultiCache] {
			e.config.Cache = anyx.IfZero(e.config.Cache, &cachex.MultiCache{})
			e.RunConfigurator(e.config.Cache)
		} else {
			var cache = &cachex.Cache{}
			e.RunConfigurator(cache, true)
			e.config.Cache = &cachex.MultiCache{cache}
		}
	}
}

// 添加自定义函数
func (e *Engine) AddCustomFunc(funcs ...func()) {
	if len(funcs) > 0 {
		e.customFuncs = append(e.customFuncs, funcs...)
	}
}

// 执行自定义函数
func (e *Engine) runCustomFunc() {
	for _, customFunc := range e.customFuncs {
		customFunc()
	}
}

// 新增自定义配置器
func (e *Engine) AddConfigurator(configurators ...confx.Configurator) {
	if len(configurators) > 0 {
		e.configurators = append(e.configurators, configurators...)
	}
}

// 运行自定义配置器
func (e *Engine) runConfigurators() {
	for _, config := range e.configurators {
		e.RunConfigurator(config)
	}
}

// 运行配置器
func (e *Engine) RunConfigurator(configurator confx.Configurator, must ...bool) {
	var ok = anyx.Default(false, must...)
	if reader := configurator.Reader(); reader != nil {
		if e.flag[EnableNacos] {
			reader.NacosGroup = e.config.Server.Name
			if err := nacosx.NewConfig(reader.NacosGroup, reader.NacosDataId, reader.Listen).Loading(configurator); err == nil {
				ok = true
			}
		} else {
			if err := marshalx.UnmarshalFromFile(e.GetConfigPath(reader.FilePath), configurator); err == nil {
				ok = true
			}
		}
	}
	if ok {
		if err := configurator.Run(); err != nil {
			log.Error(configurator.Title(), " Run Failed!")
			panic(err)
		}
		log.Info(configurator.Title(), " Run Completed!")
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
			panic("Nacos Config Client Is Uninitialized ")
		}
		var config = nacosx.NewConfig(e.config.Server.Name, dataId, listen...)
		// 加载微服务配置
		if err := config.Loading(v); err != nil {
			panic("Loading Nacos Config Failed : " + err.Error())
		}
	})
}

// 服务启动
func (e *Engine) startServer() {
	if !e.flag[NonGin] {
		defer PanicRecover()
		if e.config.Server.Debug {
			gin.SetMode(gin.DebugMode)
		}
		if e.ginEngine == nil {
			e.ginEngine = gin.New()
		}
		e.ginEngine.Use(gin.Recovery(), logx.GinRequestLog)
		e.ginEngine.Use(e.ginMiddlewares...)
		_ = e.ginEngine.SetTrustedProxies([]string{e.config.Server.Host})
		// 注册服务根路由，并执行路由注册函数
		var group = e.ginEngine.Group(e.config.Server.Prefix)
		e.InitGinLoader(group)
		var port = ":" + strconv.Itoa(e.config.Server.Port)
		log.Info("API接口请求地址: http://" + e.config.Server.Host + port)
		if err := e.ginEngine.Run(port); err != nil {
			log.Error("gin-Engine run failed ")
			panic(err)
		}
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
	if e.configDir != "" {
		return filepath.Join(e.configDir, path)
	} else {
		return path
	}
}

// 添加需要初始化的 gormx.Tabler 模型
func (e *Engine) AddTable(dst ...gormx.Tabler[any]) {
	e.AddSourceTable(constx.DefaultKey, dst...)
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
	if len(e.ginLoaders) > 0 {
		for _, loader := range e.ginLoaders {
			loader(group)
		}
	} else {
		log.Warn("engine.ginLoaders is empty ")
	}
}

// 添加需要初始化的 gormx.Tabler 模型
func (e *Engine) AddQueueTask(name string, task func()) {
	e.initQueue.AddBefore(name, task, StartServer)
}

// 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Error("server run panic: ", err)
		return
	}
	select {}
}
