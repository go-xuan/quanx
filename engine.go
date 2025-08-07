package quanx

import (
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/utilx/anyx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/filex"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/go-xuan/utilx/osx"
	"github.com/go-xuan/utilx/taskx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/serverx"
)

// Queue task name
const (
	TaskInitServer        = "init_server"         // 初始化服务
	TaskInitNacos         = "init_nacos"          // 初始化Nacos
	TaskInitConfig        = "init_config"         // 初始化配置
	TaskRunCustomFunction = "run_custom_function" // 运行自定义函数
	TaskRunServer         = "run_server"          // 运行服务
)

var engine *Engine

// GetEngine 获取当前Engine
func GetEngine() *Engine {
	if engine == nil {
		engine = NewEngine(
			EnableDebug(), // 默认启用debug
		)
	}
	return engine
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return GetEngine().config
}

// NewEngine 初始化Engine
func NewEngine(opts ...EngineOptionFunc) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			configDir:      constx.DefaultConfDir,
			customFuncs:    make([]func() error, 0),
			configurators:  make([]configx.Configurator, 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			gormTablers:    make(map[string][]interface{}),
			switches:       make(map[Option]bool),
			queue:          taskx.Queue(),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.doOptionFuncs(opts...)
		engine.queue.Add(TaskInitServer, engine.initServer)               // 1.初始化服务
		engine.queue.Add(TaskInitNacos, engine.initNacos)                 // 2.初始化Nacos
		engine.queue.Add(TaskInitConfig, engine.initConfig)               // 3.初始化配置
		engine.queue.Add(TaskRunCustomFunction, engine.runCustomFunction) // 4.运行自定义函数
		engine.queue.Add(TaskRunServer, engine.runServer)                 // 5.运行服务
	}
	return engine
}

// Engine 服务启动器
type Engine struct {
	switches       map[Option]bool          // 服务运行开关
	config         *Config                  // 服务配置数据，使用 initServer()将配置文件加载到此
	configDir      string                   // 服务配置文件夹, 使用 SetConfigDir()设置配置文件读取路径
	ginEngine      *gin.Engine              // gin框架引擎实例
	ginRouters     []func(*gin.RouterGroup) // gin路由的预加载方法，使用 AddGinRouter()添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc        // gin中间件的预加载方法，使用 AddGinRouter()添加gin中间件
	customFuncs    []func() error           // 自定义初始化函数 使用 AddCustomFunc()添加自定义函数
	configurators  []configx.Configurator   // 自定义配置器，使用 AddConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	gormTablers    map[string][]interface{} // gorm表结构对象，使用 AddTable() / AddSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
	queue          *taskx.QueueScheduler    // Engine启动时的队列任务
}

// RUN 服务运行
func (e *Engine) RUN() {
	defer PanicRecover()
	if err := e.queue.Execute(); err != nil {
		log.WithField("error", err.Error()).Error("engine run error")

	}
}

func (e *Engine) checkRunning() {
	if e.switches[running] {
		panic("engine has already running")
	}
}

// 初始化服务
func (e *Engine) initServer() error {
	e.checkRunning()
	server := &serverx.Config{}
	if e.switches[customPort] && e.config.Server != nil {
		server.Port = e.config.Server.Port
	}
	// 先设置默认值
	if err := anyx.SetDefaultValue(server); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	// 读取配置文件
	if path := e.GetConfigPath(constx.DefaultConfigFilename); filex.Exists(path) {
		if err := marshalx.Apply(path).Read(path, e.config); err != nil {
			return errorx.Wrap(err, "read file error: "+path)
		}
		if e.switches[customPort] {
			e.config.Server.Port = server.Port
		}
		if e.config.Server.Debug {
			e.switches[enableDebug] = true
		}
	} else {
		e.config.Server = server
		if err := marshalx.Apply(path).Write(path, e.config); err != nil {
			return errorx.Wrap(err, "write file error: "+path)
		}
	}
	// 更新host
	if e.config.Server.IP == "" {
		e.config.Server.IP = osx.GetLocalIP()
	}
	return nil
}

// 初始化nacos
func (e *Engine) initNacos() error {
	if e.config.Nacos != nil {
		if err := e.ExecuteConfigurator(e.config.Nacos, true); err != nil {
			return errorx.Wrap(err, "execute nacos configurator error")
		}
		// 注册nacos服务实例
		if e.config.Nacos.EnableNaming() {
			if err := nacosx.RegisterServerInstance(e.config.Server); err != nil {
				return errorx.Wrap(err, "register nacos server instance error")
			}
		}
	}
	return nil
}

// 初始化配置
func (e *Engine) initConfig() error {
	e.checkRunning()

	// 初始化日志配置器
	if err := e.initLogConfigurator(); err != nil {
		return errorx.Wrap(err, "init log configurator error")
	}

	// 初始化数据库配置器
	if err := e.initDatabaseConfigurator(); err != nil {
		return errorx.Wrap(err, "init database configurator error")
	}

	// 初始化 REDIS 缓存配置器
	if err := e.initRedisCacheConfigurator(); err != nil {
		return errorx.Wrap(err, "init cache configurator error")
	}

	// 初始化自定义配置
	for _, configurator := range e.configurators {
		if err := e.ExecuteConfigurator(configurator); err != nil {
			return errorx.Wrap(err, "init custom configurator error")
		}
	}
	return nil
}

// 运行自定义函数
func (e *Engine) runCustomFunction() error {
	e.checkRunning()
	for _, customFunc := range e.customFuncs {
		if err := customFunc(); err != nil {
			return errorx.Wrap(err, "execute custom function error")
		}
	}
	return nil
}

// 启动服务
func (e *Engine) runServer() error {
	e.checkRunning()
	if e.switches[enableDebug] {
		gin.SetMode(gin.DebugMode)
	}
	if e.ginEngine == nil {
		e.ginEngine = gin.New()
	}

	// 初始化中间件
	e.ginEngine.Use(gin.Recovery())
	if e.config.Log.Formatter == logx.FormatterJson {
		e.ginEngine.Use(ginx.JsonLogFormatter)
	} else {
		e.ginEngine.Use(ginx.DefaultLogFormatter)
	}
	e.ginEngine.Use(e.ginMiddlewares...)

	server := e.config.Server
	if err := e.ginEngine.SetTrustedProxies([]string{server.IP}); err != nil {
		return errorx.Wrap(err, "set trusted proxies error")
	}
	// 注册服务根路由
	group := e.ginEngine.Group(server.ApiPrefix())
	e.initGinRouter(group)

	// 启动服务
	e.switches[running] = true
	log.Infof(`API接口请求地址: %s`, server.ApiDomain())
	if err := e.ginEngine.Run(":" + strconv.Itoa(server.Port)); err != nil {
		return errorx.Wrap(err, "run server error")
	}
	return nil
}

// 初始化日志
func (e *Engine) initLogConfigurator() error {
	config := e.config.Log
	if config == nil {
		config = logx.GetConfig()
	}
	if config.Name == "" || config.Name == "app" {
		config.Name = e.config.Server.Name
	}
	if err := e.ExecuteConfigurator(config, true); err != nil {
		return errorx.Wrap(err, "execute log configurator error")
	}
	e.config.Log = config
	return nil
}

// 初始化数据库
func (e *Engine) initDatabaseConfigurator() error {
	// 读取数据库配置并初始化
	database := anyx.IfZero(e.config.Database, &gormx.Configs{})
	if err := e.ExecuteConfigurator(database); err != nil {
		e.config.Database = database
	}
	if !gormx.Initialized() {
		var one = &gormx.Config{}
		if err := e.ExecuteConfigurator(one); err == nil {
			e.config.Database = &gormx.Configs{one}
		}
	}

	// 初始化表结构
	if gormx.Initialized() {
		for _, source := range gormx.Sources() {
			if tablers, ok := e.gormTablers[source]; ok {
				if err := gormx.InitTable(source, tablers...); err != nil {
					return errorx.Wrap(err, "init table struct and data failed")
				}
			}
		}
	}

	return nil
}

// 初始化缓存
func (e *Engine) initRedisCacheConfigurator() error {
	// 读取redis配置并初始化
	redis := anyx.IfZero(e.config.Redis, &redisx.Configs{})
	if err := e.ExecuteConfigurator(redis); err == nil {
		e.config.Redis = redis
	}
	if !redisx.Initialized() {
		var one = &redisx.Config{}
		if err := e.ExecuteConfigurator(one); err == nil {
			e.config.Redis = &redisx.Configs{one}
		}
	}

	// 初始化缓存
	if redisx.Initialized() {
		cache := anyx.IfZero(e.config.Cache, &cachex.Configs{})
		if err := e.ExecuteConfigurator(cache); err == nil {
			e.config.Cache = cache
		}
		if !cachex.Initialized() {
			var one = &cachex.Config{}
			if err := e.ExecuteConfigurator(one); err == nil {
				e.config.Cache = &cachex.Configs{one}
			}
		}
	}
	return nil
}

// initGinRouter 初始化gin路由
func (e *Engine) initGinRouter(group *gin.RouterGroup) {
	if len(e.ginRouters) > 0 {
		for _, router := range e.ginRouters {
			router(group)
		}
	} else {
		log.Warn("gin router is empty")
	}
}

// AddCustomFunc 添加自定义函数
func (e *Engine) AddCustomFunc(funcs ...func() error) {
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

func (e *Engine) ReadConfigurator(configurator configx.Configurator) (string, error) {
	e.checkRunning()
	// 获取可用的reader，优先级：nacos > file > env > tag
	if reader := func() configx.Reader {
		if reader := configx.CheckReader(configurator, configx.FromNacos); reader != nil {
			reader.Anchor(e.config.Server.NacosGroup())
			if err := reader.Read(configurator); err == nil {
				return reader
			}
		}
		if reader := configx.CheckReader(configurator, configx.FromFile); reader != nil {
			reader.Anchor(e.configDir)
			if err := reader.Read(configurator); err == nil {
				return reader
			}
		}
		if reader := configx.CheckReader(configurator, configx.FromEnv); reader != nil {
			if err := reader.Read(configurator); err == nil {
				return reader
			}
		}
		if reader := configx.CheckReader(configurator, configx.FromTag); reader != nil {
			reader.Anchor("default")
			if err := reader.Read(configurator); err == nil {
				return reader
			}
		}
		return nil
	}(); reader != nil {
		return reader.Location(), nil
	}
	return "", errorx.New("configurator reader is nil")
}

// ExecuteConfigurator 执行配置器
// configurator: 配置器
// must: 是否必须执行配置器，默认false
func (e *Engine) ExecuteConfigurator(configurator configx.Configurator, must ...bool) error {
	var logger = log.WithField("location", "none")
	if configurator.NeedRead() {
		if location, err := e.ReadConfigurator(configurator); err != nil && (len(must) == 0 || !must[0]) {
			// 如果是非必须执行的配置器，读取配置失败则直接返回
			return errorx.Wrap(err, "configurator read error")
		} else {
			logger = logger.WithField("location", location)
		}
	}
	// 执行配置器
	if err := configurator.Execute(); err != nil {
		logger.WithField("error", err.Error()).Error("configurator execute error")
		return errorx.Wrap(err, "configurator execute error")
	}
	if e.switches[enableDebug] {
		logger = logger.WithField("type", reflect.TypeOf(configurator).String())
	}
	logger.Info("configurator execute success")
	return nil
}

// ReadLocalConfig 读取本地配置项（立即执行）
func (e *Engine) ReadLocalConfig(config any, path string) {
	if err := marshalx.Apply(path).Read(path, config); err != nil {
		log.WithField("path", path).WithField("error", err.Error()).Error("read local config error")
		panic(errorx.Wrap(err, "read local config error"))
	}
}

// ReadNacosConfig 读取nacos配置（以自定义函数的形式延迟执行，需确保nacos已经提前初始化）
func (e *Engine) ReadNacosConfig(config any, dataId string, listen ...bool) {
	e.AddCustomFunc(func() error {
		var reader = &nacosx.Reader{
			DataId: dataId,
			Group:  e.config.Server.NacosGroup(),
			Listen: len(listen) > 0 && listen[0],
		}
		if err := reader.Read(config); err != nil {
			return errorx.Wrap(err, "read nacos config error")
		}
		return nil
	})
}

// doOptionFuncs 加载配置函数
func (e *Engine) doOptionFuncs(funcs ...EngineOptionFunc) {
	e.checkRunning()
	if len(funcs) > 0 {
		for _, f := range funcs {
			f(e)
		}
	}
}

// SetConfigDir 设置配置文件路径
func (e *Engine) SetConfigDir(dir string) {
	e.checkRunning()
	e.configDir = dir
}

// GetConfigPath 获取配置文件
func (e *Engine) GetConfigPath(filename string) string {
	if dir := e.configDir; dir != "" {
		return filepath.Join(dir, filename)
	} else {
		return filename
	}
}

// AddTable 添加表结构（默认数据源）
func (e *Engine) AddTable(tablers ...interface{}) {
	e.AddSourceTable("default", tablers...)
}

// AddSourceTable 添加表结构（指定数据源）
func (e *Engine) AddSourceTable(source string, tablers ...interface{}) {
	e.checkRunning()
	if len(tablers) > 0 {
		e.gormTablers[source] = append(e.gormTablers[source], tablers...)
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

// AddTaskBefore 前插队添加任务
func (e *Engine) AddTaskBefore(base, name string, task func() error) {
	e.checkRunning()
	e.queue.AddBefore(base, name, task)
}

// AddTaskAfter 后插队添加任务
func (e *Engine) AddTaskAfter(base, name string, task func() error) {
	e.checkRunning()
	e.queue.AddAfter(base, name, task)
}

// PanicRecover 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		log.Error("engine run panic: ", err)
		return
	}
}
