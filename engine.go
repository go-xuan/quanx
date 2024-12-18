package quanx

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/core/cachex"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/ginx"
	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/core/nacosx"
	"github.com/go-xuan/quanx/core/redisx"
	"github.com/go-xuan/quanx/net/ipx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/os/filex"
	"github.com/go-xuan/quanx/os/fmtx"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/os/syncx"
	"github.com/go-xuan/quanx/os/taskx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/types/stringx"
	"github.com/go-xuan/quanx/utils/marshalx"
)

// Mode 服务启动模式
type Mode uint

const (
	Debug         Mode = iota // debug模式
	EnableNacos               // 启用nacos
	MultiDatabase             // 开启多数据源
	MultiRedis                // 开启多redis源
	MultiCache                // 开启多缓存源
	EnableQueue               // 使用队列任务启动
	Running                   // 正在运行中
)

// Init Queue Function Department
const (
	LoadingConfig        = "loading_config"        // 加载配置文件
	InitServerBasic      = "init_server_basic"     // 初始化服务基础组件（log/nacos/gorm/redis/cache）
	ExecuteConfigurators = "execute_configurators" // 运行自定义配置器
	RunCustomFuncs       = "run_custom_funcs"      // 运行自定义函数
	StartServer          = "start_server"          // 服务启动
)

var engine *Engine

// Engine 服务配置器
type Engine struct {
	mode           map[Mode]bool             // 服务运行模式
	config         *Config                   // 服务配置数据，使用 loadingAppConfig()将配置文件加载到此
	configDir      string                    // 服务配置文件夹, 使用 SetConfigDir()设置配置文件读取路径
	ginEngine      *gin.Engine               // gin框架引擎实例
	ginRouters     []func(*gin.RouterGroup)  // gin路由的预加载方法，使用 AddGinRouter()添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc         // gin中间件的预加载方法，使用 AddGinRouter()添加gin中间件
	customFuncs    []func()                  // 自定义初始化函数 使用 AddCustomFunc()添加自定义函数
	configurators  []configx.Configurator    // 配置器，使用 AddConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	gormTables     map[string][]gormx.Tabler // gorm表结构对象，使用 AddTable() / AddSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
	queue          *taskx.QueueScheduler     // Engine启动时的队列任务
}

// GetEngine 获取当前Engine
func GetEngine() *Engine {
	if engine == nil {
		engine = DefaultEngine()
	}
	return engine
}

// DefaultEngine 默认Engine
func DefaultEngine() *Engine {
	return NewEngine(
		Debug,
	)
}

// NewEngine 初始化Engine
func NewEngine(modes ...Mode) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			configDir:      constx.DefaultConfDir,
			customFuncs:    make([]func(), 0),
			configurators:  make([]configx.Configurator, 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			gormTables:     make(map[string][]gormx.Tabler),
			mode:           make(map[Mode]bool),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.SetMode(modes...)
	}
	// 设置默认日志输出
	log.SetOutput(logx.DefaultWriter())
	log.SetFormatter(logx.DefaultFormatter())
	// 设置服务启动队列
	engine.enableQueue()
	return engine
}

// RUN 服务运行
func (e *Engine) RUN() {
	if engine.mode[EnableQueue] { // 任务队列方式启动
		engine.queue.Execute()
	} else { // 默认方式启动
		syncx.OnceDo(e.loadingAppConfig)     // 1.加载服务配置文件
		syncx.OnceDo(e.initAppBasic)         // 2.初始化服务基础组件（log/nacos/gorm/redis/cache）
		syncx.OnceDo(e.executeConfigurators) // 3.运行自定义配置器
		syncx.OnceDo(e.runCustomFuncs)       // 4.运行自定义函数
		syncx.OnceDo(e.startServer)          // 5.服务启动
	}
	e.mode[Running] = true
}

func (e *Engine) checkRunning() {
	if engine.mode[Running] {
		panic("engine has already running")
	}
}

// 是否启用队列
func (e *Engine) enableQueue() {
	if engine.mode[EnableQueue] && e.queue == nil {
		queue := taskx.Queue()
		queue.Add(LoadingConfig, engine.loadingAppConfig)            // 1.加载服务配置文件
		queue.Add(InitServerBasic, engine.initAppBasic)              // 2.初始化服务基础组件（log/nacos/gorm/redis/cache）
		queue.Add(ExecuteConfigurators, engine.executeConfigurators) // 3.运行自定义配置器
		queue.Add(RunCustomFuncs, engine.runCustomFuncs)             // 4.运行自定义函数
		queue.Add(StartServer, engine.startServer)                   // 5.服务启动
		engine.queue = queue
	}
}

// 加载服务配置文件
func (e *Engine) loadingAppConfig() {
	e.checkRunning()
	var server = &Server{}
	// 先设置默认值
	if err := anyx.SetDefaultValue(server); err != nil {
		panic(errorx.Wrap(err, "set default value error"))
	}
	var config = e.config
	config.Server = server
	// 读取配置文件
	if path := e.GetConfigPath(constx.DefaultConfigFilename); filex.Exists(path) {
		if err := marshalx.UnmarshalFromFile(path, config); err != nil {
			panic(errorx.Wrap(err, "unmarshal file failed: "+path))
		}
	} else if err := marshalx.WriteYaml(path, config); err != nil {
		panic(errorx.Wrap(err, "set default value error"))
	}
	if config.Server.Host == "" {
		config.Server.Host = ipx.GetLocalIP()
	}
	// 从nacos加载配置
	if e.mode[EnableNacos] && config.Nacos != nil {
		e.ExecuteConfigurator(config.Nacos, true)
		if config.Nacos.EnableNaming() {
			// 注册nacos服务Nacos
			if err := nacosx.Register(config.Server.Instance()); err != nil {
				panic(errorx.Wrap(err, "nacos register error"))
			}
		}
	}
	e.config = config
}

// 初始化服务基础组件（log/nacos/gorm/redis）
func (e *Engine) initAppBasic() {
	e.checkRunning()
	// 初始化日志
	var serverName = stringx.IfZero(e.config.Server.Name, "app")
	e.ExecuteConfigurator(anyx.IfZero(e.config.Log, &logx.Config{FileName: serverName + ".log"}), true)

	// 初始化数据库连接
	if e.config.Database != nil {
		e.ExecuteConfigurator(e.config.Database)
	} else if e.mode[MultiDatabase] {
		var multi = &gormx.MultiConfig{}
		e.ExecuteConfigurator(multi)
		e.config.Database = multi
	} else {
		var db = &gormx.Config{}
		e.ExecuteConfigurator(db)
		e.config.Database = &gormx.MultiConfig{db}
	}

	// 初始化表结构
	if gormx.Initialized() {
		for _, source := range gormx.Sources() {
			if dst, ok := e.gormTables[source]; ok {
				if err := gormx.InitTable(source, dst...); err != nil {
					panic(errorx.Wrap(err, "init table struct and data failed"))
				}
			}
		}
	}

	// 初始化redis连接

	if e.config.Redis != nil {
		e.ExecuteConfigurator(e.config.Redis)
	} else if e.mode[MultiRedis] {
		var multi = &redisx.MultiConfig{}
		e.ExecuteConfigurator(multi)
		e.config.Redis = multi
	} else {
		var db = &redisx.Config{}
		e.ExecuteConfigurator(db)
		e.config.Redis = &redisx.MultiConfig{db}
	}

	// 初始化缓存
	if redisx.Initialized() {
		if e.mode[MultiCache] {
			e.config.Cache = anyx.IfZero(e.config.Cache, &cachex.MultiConfig{})
			e.ExecuteConfigurator(e.config.Cache)
		} else {
			var cache = &cachex.Config{}
			e.ExecuteConfigurator(cache, true)
			e.config.Cache = &cachex.MultiConfig{cache}
		}
	}
}

// 运行自定义配置器
func (e *Engine) executeConfigurators() {
	e.checkRunning()
	for _, configurator := range e.configurators {
		e.ExecuteConfigurator(configurator)
	}
}

// 运行自定义函数
func (e *Engine) runCustomFuncs() {
	e.checkRunning()
	for _, customFunc := range e.customFuncs {
		customFunc()
	}
}

// 服务启动
func (e *Engine) startServer() {
	e.checkRunning()
	defer PanicRecover()
	if e.mode[Debug] {
		gin.SetMode(gin.DebugMode)
	}
	if e.ginEngine == nil {
		e.ginEngine = gin.New()
	}
	e.ginEngine.Use(gin.Recovery(), ginx.RequestLogFmt)
	e.ginEngine.Use(e.ginMiddlewares...)
	_ = e.ginEngine.SetTrustedProxies([]string{e.config.Server.Host})
	// 注册服务根路由
	group := e.ginEngine.Group(e.config.Server.ApiPrefix())
	e.InitGinLoader(group)
	log.Info("API接口请求地址: " + e.config.Server.ApiHost())
	if err := e.ginEngine.Run(":" + strconv.Itoa(e.config.Server.Port)); err != nil {
		panic(errorx.Wrap(err, "gin run failed"))
	}
}

// AddCustomFunc 添加自定义函数
func (e *Engine) AddCustomFunc(funcs ...func()) {
	e.checkRunning()
	if len(funcs) > 0 {
		e.customFuncs = append(e.customFuncs, funcs...)
	}
}

// AddConfigurator 新增自定义配置器
func (e *Engine) AddConfigurator(configurators ...configx.Configurator) {
	e.checkRunning()
	if len(configurators) > 0 {
		e.configurators = append(e.configurators, configurators...)
	}
}

// ExecuteConfigurator 运行配置器
func (e *Engine) ExecuteConfigurator(configurator configx.Configurator, must ...bool) {
	e.checkRunning()
	var mustRun = anyx.Default(false, must...)
	if reader := configurator.Reader(); reader != nil {
		if e.mode[EnableNacos] {
			reader.NacosGroup = e.config.Server.Name
			if err := nacosx.This().ScanConfig(configurator, reader.NacosGroup, reader.NacosDataId, reader.Listen); err == nil {
				mustRun = true
			}
		} else {
			if err := marshalx.UnmarshalFromFile(e.GetConfigPath(reader.FilePath), configurator); err == nil {
				mustRun = true
			}
		}
	}
	if mustRun {
		var logger = log.WithField("configurator", configurator.ID())
		if err := configx.Execute(configurator); err != nil {
			if e.mode[Debug] {
				logger.Info("configurator data: ", configurator.Format())
			}
			logger.Error(fmtx.Red.String("configurator execute failed!"))
			log.Panic(err)
		}
		if e.mode[Debug] {
			logger.Info("configurator data: ", configurator.Format())
		}
		logger.Info(fmtx.Green.String("configurator execute success!"))
	}
}

// LoadingLocalConfig 加载本地配置项（立即加载）
func (e *Engine) LoadingLocalConfig(v any, path string) {
	if err := marshalx.UnmarshalFromFile(path, v); err != nil {
		panic(errorx.Wrap(err, "unmarshal config file failed"))
	}
}

// LoadingNacosConfig 加载nacos配置（以自定义函数的形式延迟加载）
func (e *Engine) LoadingNacosConfig(v any, dataId string, listen ...bool) {
	e.AddCustomFunc(func() {
		if err := nacosx.This().ScanConfig(v, e.config.Server.Name, dataId, listen...); err != nil {
			panic(errorx.Wrap(err, "scan nacos config failed"))
		}
	})
}

// SetMode 设置模式
func (e *Engine) SetMode(flags ...Mode) {
	e.checkRunning()
	if len(flags) > 0 {
		for _, flag := range flags {
			e.mode[flag] = true
		}
	}
}

// SetConfigDir 设置配置文件
func (e *Engine) SetConfigDir(dir string) {
	e.checkRunning()
	e.configDir = dir
}

// GetConfigPath 设置配置文件
func (e *Engine) GetConfigPath(filename string) string {
	if dir := e.configDir; dir != "" {
		return filepath.Join(dir, filename)
	} else {
		return filename
	}
}

// AddTable 添加需要初始化的 gormx.Tabler 模型
func (e *Engine) AddTable(dst ...gormx.Tabler) {
	e.AddSourceTable(constx.DefaultSource, dst...)
}

// AddSourceTable 添加需要某个数据源的gormx.Table模型
func (e *Engine) AddSourceTable(source string, dst ...gormx.Tabler) {
	e.checkRunning()
	if len(dst) > 0 {
		e.gormTables[source] = append(e.gormTables[source], dst...)
	}
}

// AddGinMiddleware 添加gin中间件
func (e *Engine) AddGinMiddleware(middleware ...gin.HandlerFunc) {
	e.checkRunning()
	if len(middleware) > 0 {
		e.ginMiddlewares = append(e.ginMiddlewares, middleware...)
	}
}

// AddGinRouter 添加gin的路由加载函数
func (e *Engine) AddGinRouter(router ...func(*gin.RouterGroup)) {
	e.checkRunning()
	if len(router) > 0 {
		e.ginRouters = append(e.ginRouters, router...)
	}
}

// InitGinLoader 执行gin的路由加载函数
func (e *Engine) InitGinLoader(group *gin.RouterGroup) {
	if len(e.ginRouters) > 0 {
		for _, loader := range e.ginRouters {
			loader(group)
		}
	} else {
		log.Warn("gin router is empty")
	}
}

// AddQueueTask 使用后，会自动启用队列方式启动服务，且当前添加的任务会在 StartServer 之前执行
func (e *Engine) AddQueueTask(name string, task func()) {
	e.checkRunning()
	if name == "" {
		log.Error(`add queue task failed, cause: the task name is required`)
	} else {
		e.mode[EnableQueue] = true
		e.enableQueue()
		e.queue.AddBefore(name, task, StartServer)
		log.Info(`add queue task successfully, task name:`, name)
	}
	return
}

// PanicRecover 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Error("server run panic: ", err)
		return
	}
	select {}
}
