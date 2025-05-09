package quanx

import (
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/base/errorx"
	"github.com/go-xuan/quanx/base/filex"
	"github.com/go-xuan/quanx/base/osx"
	"github.com/go-xuan/quanx/base/taskx"
	"github.com/go-xuan/quanx/common/constx"
	"github.com/go-xuan/quanx/extra/cachex"
	"github.com/go-xuan/quanx/extra/configx"
	"github.com/go-xuan/quanx/extra/ginx"
	"github.com/go-xuan/quanx/extra/gormx"
	"github.com/go-xuan/quanx/extra/logx"
	"github.com/go-xuan/quanx/extra/nacosx"
	"github.com/go-xuan/quanx/extra/redisx"
	"github.com/go-xuan/quanx/types/anyx"
	"github.com/go-xuan/quanx/utils/marshalx"
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

// GetEngine 获取当前Engine
func GetEngine() *Engine {
	if engine == nil {
		engine = NewEngine(
			EnableDebug(), // 默认启用debug
		)
	}
	return engine
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

// RUN 服务运行
func (e *Engine) RUN() {
	defer PanicRecover()
	if err := e.queue.Execute(); err != nil {
		panic(err)
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
	server := &Server{}
	if e.switches[customPort] && e.config.Server != nil {
		server.Port = e.config.Server.Port
	}
	// 先设置默认值
	if err := anyx.SetDefaultValue(server); err != nil {
		return errorx.Wrap(err, "set default value error")
	}
	// 读取公共配置文件
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
			return errorx.Wrap(err, "Write file error: "+path)
		}
	}
	// 更新host
	if e.config.Server.Host == "" {
		e.config.Server.Host = osx.GetLocalIP()
	}
	return nil
}

// 初始化nacos
func (e *Engine) initNacos() error {
	if e.config.Nacos != nil {
		if err := e.ExecuteConfigurator(e.config.Nacos, true); err != nil {
			return errorx.Wrap(err, "execute nacos configurator error")
		}
		if e.config.Nacos.EnableNaming() { // 注册nacos服务
			if err := nacosx.Register(e.config.Server.Instance()); err != nil {
				return errorx.Wrap(err, "register nacos error")
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

	// 初始化 Redis 缓存配置器
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

	host := e.config.Server.Host
	if err := e.ginEngine.SetTrustedProxies([]string{host}); err != nil {
		return errorx.Wrap(err, "set trusted proxies error")
	}
	// 注册服务根路由
	group := e.ginEngine.Group(e.config.Server.ApiPrefix())
	e.initGinRouter(group)

	// 获取服务端口
	port := strconv.Itoa(e.config.Server.Port)
	// 启动服务
	e.switches[running] = true
	logrus.Infof(`API接口请求地址: http://%s:%s`, host, port)
	if err := e.ginEngine.Run(":" + port); err != nil {
		return errorx.Wrap(err, "run server error")
	}
	return nil
}

// 初始化日志
func (e *Engine) initLogConfigurator() error {
	if e.config.Log == nil {
		e.config.Log = logx.GetConfig()
		return nil
	}
	if e.config.Log.Name == "" {
		e.config.Log.Name = e.config.Server.Name
	}
	if err := e.ExecuteConfigurator(e.config.Log, true); err != nil {
		return errorx.Wrap(err, "execute log configurator error")
	}
	return nil
}

// 初始化数据库
func (e *Engine) initDatabaseConfigurator() error {
	// 读取数据库配置并初始化
	database := anyx.IfZero(e.config.Database, &gormx.MultiConfig{})
	if err := e.ExecuteConfigurator(database); err != nil {
		e.config.Database = database
	}
	if !gormx.Initialized() {
		var one = &gormx.Config{}
		if err := e.ExecuteConfigurator(one); err == nil {
			e.config.Database = &gormx.MultiConfig{one}
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
	redis := anyx.IfZero(e.config.Redis, &redisx.MultiConfig{})
	if err := e.ExecuteConfigurator(redis); err == nil {
		e.config.Redis = redis
	}
	if !redisx.Initialized() {
		var one = &redisx.Config{}
		if err := e.ExecuteConfigurator(one); err == nil {
			e.config.Redis = &redisx.MultiConfig{one}
		}
	}

	// 初始化缓存
	if redisx.Initialized() {
		cache := anyx.IfZero(e.config.Cache, &cachex.MultiConfig{})
		if err := e.ExecuteConfigurator(cache); err == nil {
			e.config.Cache = cache
		}
		if !cachex.Initialized() {
			var one = &cachex.Config{}
			if err := e.ExecuteConfigurator(one); err == nil {
				e.config.Cache = &cachex.MultiConfig{one}
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
		logrus.Warn("gin router is empty")
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
	// 获取可用的reader
	reader := func() configx.Reader {
		if reader := configurator.Reader(configx.FormNacos); reader != nil {
			reader.Ready(e.config.Server.Name)
			if err := reader.Check(configurator); err == nil {
				return reader
			}
		}
		if reader := configurator.Reader(configx.FromLocal); reader != nil {
			reader.Ready(e.configDir)
			if err := reader.Check(configurator); err == nil {
				return reader
			}
		}
		if reader := configurator.Reader(configx.FromEnv); reader != nil {
			reader.Ready()
			if err := reader.Check(configurator); err == nil {
				return reader
			}
		}
		return nil
	}()
	if reader == nil {
		return configx.FromDefault, errorx.New("get config reader error")
	}
	return reader.Location(), reader.Read(configurator)
}

// ExecuteConfigurator 执行配置器（立即执行）
func (e *Engine) ExecuteConfigurator(configurator configx.Configurator, must ...bool) error {
	location, err := e.ReadConfigurator(configurator)
	logger := logrus.WithField("location", location)
	if err != nil && (len(must) == 0 || !must[0]) {
		// 配置读取失败并且非必须执行，则直接返回
		return errorx.Wrap(err, "configurator read error")
	}
	if err = configurator.Execute(); err != nil {
		logger.Error("configurator execute error")
		return errorx.Wrap(err, "configurator execute error")
	}
	if e.switches[enableDebug] {
		logger = logger.WithField("info", configurator.Info())
	}
	logger.Info("configurator execute success")
	return nil
}

// ReadLocalConfig 读取本地配置项（立即执行）
func (e *Engine) ReadLocalConfig(config any, path string) {
	if err := marshalx.Apply(path).Read(path, config); err != nil {
		logrus.WithField("path", path).WithError(err).Error("read local config error")
		panic(errorx.Wrap(err, "read local config error"))
	}
}

// ReadNacosConfig 读取nacos配置（以自定义函数的形式延迟执行，需确保nacos已经提前初始化）
func (e *Engine) ReadNacosConfig(config any, dataId string, listen ...bool) {
	e.AddCustomFunc(func() error {
		if !nacosx.Initialized() {
			var err = errorx.New("nacos not initialized")
			logrus.WithField("dataId", dataId).WithError(err).Error("read nacos config error")
			return err
		}
		if err := nacosx.ReadConfig(config, e.config.Server.Name, dataId, listen...); err != nil {
			return errorx.Wrap(err, "read nacos config error")
		}
		return nil
	})
}

// doOptionFuncs 启动项函数
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
	e.AddSourceTable(constx.DefaultSource, tablers...)
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
func (e *Engine) AddTaskBefore(name, before string, task func() error) {
	e.checkRunning()
	e.queue.AddBefore(name, before, task)
}

// AddTaskAfter 后插队添加任务
func (e *Engine) AddTaskAfter(name, after string, task func() error) {
	e.checkRunning()
	e.queue.AddAfter(name, after, task)
}

// PanicRecover 服务保活
func PanicRecover() {
	if err := recover(); err != nil {
		logrus.Error("server run panic: ", err)
		return
	}
}
