package quanx

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/go-xuan/utilx/taskx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/ginx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
)

// queue task name
const (
	StepInitConfig  = "init_config"  // 初始化配置
	StepRunExecutes = "run_executes" // 运行自定义函数
	StepRunServer   = "run_server"   // 运行服务
	StepRunning     = "running"
)

var engine *Engine

// GetEngine 获取当前Engine
func GetEngine() *Engine {
	if engine == nil {
		engine = NewEngine(
			SetPort(constx.DefaultPort),
			Debug(),
		)
	}
	return engine
}

// NewEngine 初始化Engine
func NewEngine(opts ...EngineOption) *Engine {
	if engine == nil {
		engine = &Engine{
			config:         &Config{},
			executes:       make([]taskx.Execute, 0),
			configurators:  make([]configx.Configurator, 0),
			ginMiddlewares: make([]gin.HandlerFunc, 0),
			tablers:        make(map[string][]interface{}),
			flags:          make(map[string]bool),
			queue:          taskx.NewQueueScheduler("engine"),
		}
		gin.SetMode(gin.ReleaseMode)
		engine.DoOption(opts...)
		engine.queue.Add(StepInitConfig, engine.initConfig)   // 1.初始化配置
		engine.queue.Add(StepRunExecutes, engine.runExecutes) // 2.运行自定义函数
		engine.queue.Add(StepRunServer, engine.runServer)     // 3.运行服务
	}
	return engine
}

// Engine 服务启动器
type Engine struct {
	config         *Config                  // 服务配置，使用 initServer()加载配置
	ginEngine      *gin.Engine              // gin框架引擎实例
	ginRouters     []func(*gin.RouterGroup) // gin路由的预加载方法，使用 addGinRouter()添加自行实现的路由注册方法
	ginMiddlewares []gin.HandlerFunc        // gin中间件的预加载方法，使用 addGinRouter()添加gin中间件
	executes       []taskx.Execute          // 自定义初始化函数 使用 addExecute()添加自定义函数
	tablers        map[string][]interface{} // gorm表结构对象，使用 addTable() / addSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
	configurators  []configx.Configurator   // 自定义配置器，使用 addConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	flags          map[string]bool          // 服务运行标识
	queue          *taskx.QueueScheduler    // Engine启动步骤队列
}

// GetConfig 获取当前配置
func (e *Engine) GetConfig() *Config {
	return e.config
}

// DoOption 加载配置选项
func (e *Engine) DoOption(options ...EngineOption) {
	e.runningCheck()
	for _, option := range options {
		option(e)
	}
}

// RUN 服务运行
func (e *Engine) RUN(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("engine run panic: ", err)
			return
		}
	}()

	// 以队列形式依次执行启动步骤
	if err := e.queue.Execute(ctx); err != nil {
		log.WithField("error", err.Error()).Error("engine run error")
	}
}

// 检查服务运行状态
func (e *Engine) runningCheck() {
	if e.flags[StepRunning] {
		panic("server has already running")
	}
}

// 初始化配置
func (e *Engine) initConfig(ctx context.Context) error {
	e.runningCheck()
	e.flags[StepInitConfig] = true

	// 初始化配置
	if err := e.config.Init(constx.GetDefaultConfigPath()); err != nil {
		return errorx.Wrap(err, "init config error")
	}

	// 初始化数据库表结构
	if err := e.initTables(); err != nil {
		return errorx.Wrap(err, "init tables error")
	}

	// 运行自定义配置器
	for _, configurator := range e.configurators {
		if err := configx.ConfiguratorReadAndExecute(configurator); err != nil {
			return errorx.Wrap(err, "init configurators error")
		}
	}

	return nil
}

// 初始化数据库表
func (e *Engine) initTables() error {
	if gormx.Initialized() {
		for _, source := range gormx.Sources() {
			if tablers, ok := e.tablers[source]; ok {
				if err := gormx.InitTable(source, tablers...); err != nil {
					return errorx.Wrap(err, "init table struct and data failed")
				}
			}
		}
	}

	return nil
}

// 运行自定义函数
func (e *Engine) runExecutes(ctx context.Context) error {
	e.runningCheck()
	e.flags[StepRunExecutes] = true
	for _, execute := range e.executes {
		if err := execute(ctx); err != nil {
			return errorx.Wrap(err, "execute error")
		}
	}
	return nil
}

// 启动服务
func (e *Engine) runServer(ctx context.Context) error {
	e.runningCheck()
	e.flags[StepRunServer] = true

	if e.config.Server.Debug {
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
	if err := e.ginEngine.SetTrustedProxies([]string{server.GetIP()}); err != nil {
		return errorx.Wrap(err, "set trusted proxies error")
	}
	// 注册服务根路由
	group := e.ginEngine.Group(server.ApiPrefix())
	e.initGinRouter(group)

	// 启动服务
	e.flags[StepRunning] = true
	log.Infof(`service address: %s`, server.GetAddress())
	if err := e.ginEngine.Run(":" + strconv.Itoa(server.Port)); err != nil {
		return errorx.Wrap(err, "run server error")
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

// 新增自定义配置器
func (e *Engine) addConfigurator(configurators ...configx.Configurator) {
	e.runningCheck()
	if len(configurators) > 0 {
		e.configurators = append(e.configurators, configurators...)
	}
}

// 添加自定义函数
func (e *Engine) addExecute(executes ...taskx.Execute) {
	e.runningCheck()
	e.executes = append(e.executes, executes...)
}

// 添加表结构（默认数据源）
func (e *Engine) addTable(tablers ...interface{}) {
	e.addSourceTable("default", tablers...)
}

// 添加表结构（指定数据源）
func (e *Engine) addSourceTable(source string, tablers ...interface{}) {
	e.runningCheck()
	if len(tablers) > 0 {
		e.tablers[source] = append(e.tablers[source], tablers...)
	}
}

// 添加gin的路由加载函数
func (e *Engine) addGinRouter(router ...func(*gin.RouterGroup)) {
	e.runningCheck()
	if len(router) > 0 {
		e.ginRouters = append(e.ginRouters, router...)
	}
}

// 添加gin中间件
func (e *Engine) addGinMiddleware(middleware ...gin.HandlerFunc) {
	e.runningCheck()
	if len(middleware) > 0 {
		e.ginMiddlewares = append(e.ginMiddlewares, middleware...)
	}
}

// 前插队添加任务
func (e *Engine) addTaskBefore(base, name string, task func(context.Context) error) {
	e.runningCheck()
	e.queue.AddBefore(base, name, task)
}

// 后插队添加任务
func (e *Engine) addTaskAfter(base, name string, task func(context.Context) error) {
	e.runningCheck()
	e.queue.AddAfter(base, name, task)
}

// ReadLocalConfig 读取本地配置项（立即执行）
func (e *Engine) readLocalConfig(config any, path string) {
	if err := marshalx.Apply(path).Read(path, config); err != nil {
		log.WithField("path", path).WithField("error", err.Error()).
			Error("read local config error")
		panic(errorx.Wrap(err, "read local config error"))
	}
}

// 读取nacos配置（需保证nacos已经提前初始化，所以以自定义函数的形式延迟执行）
func (e *Engine) readNacosConfig(config any, dataId string, listen ...bool) {
	e.addExecute(func(ctx context.Context) error {
		if err := configx.ReadWithReader(config, &nacosx.Reader{
			DataId: dataId,
			Listen: len(listen) > 0 && listen[0],
		}); err != nil {
			log.WithField("dataId", dataId).WithField("error", err.Error()).
				Error("read nacos config error")
			return errorx.Wrap(err, "read nacos config error")
		}
		return nil
	})
}
