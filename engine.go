package quanx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-xuan/quanx/serverx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/marshalx"
	"github.com/go-xuan/utilx/taskx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/nacosx"
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
			InitServer(&serverx.Config{Debug: true}),
		)
	}
	return engine
}

// NewEngine 初始化Engine
func NewEngine(opts ...EngineOption) *Engine {
	if engine == nil {
		engine = &Engine{
			config:        &Config{},
			executes:      make([]taskx.Execute, 0),
			servers:       make([]serverx.Server, 0),
			configurators: make([]configx.Configurator, 0),
			tablers:       make(map[string][]interface{}),
			flags:         make(map[string]bool),
		}
		// 初始化队列
		engine.initQueue()
		// 加载配置选项
		engine.LoadOption(opts...)
	}
	return engine
}

// Engine 服务启动器
type Engine struct {
	config        *Config                  // 服务配置
	servers       []serverx.Server         // http/grpc服务
	executes      []taskx.Execute          // 自定义初始化函数 使用 addExecute()添加自定义函数
	configurators []configx.Configurator   // 自定义配置器，使用 addConfigurator()添加配置器对象，被添加对象必须为指针类型，且需要实现 configx.Configurator 接口
	tablers       map[string][]interface{} // gorm表结构对象，使用 addTable() / addSourceTable() 添加至表结构初始化任务列表，需要实现 gormx.Tabler 接口
	flags         map[string]bool          // 服务运行标识
	queue         *taskx.QueueScheduler    // Engine启动步骤队列
}

// LoadOption 加载配置选项
func (e *Engine) LoadOption(options ...EngineOption) {
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

// InitQueue 初始化队列
func (e *Engine) initQueue() {
	queue := taskx.NewQueueScheduler("engine")
	queue.Add(StepInitConfig, engine.initConfig)   // 1.初始化配置
	queue.Add(StepRunExecutes, engine.runExecutes) // 2.运行自定义函数
	queue.Add(StepRunServer, engine.runServer)     // 3.运行服务
	e.queue = queue
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

	// 启动服务
	for _, server := range e.servers {
		if err := server.Run(e.config.Server); err != nil {
			return errorx.Wrap(err, "server run error")
		}
	}

	e.flags[StepRunning] = true

	quit := make(chan os.Signal, 1)
	// 接收 SIGINT（Ctrl+C）、SIGTERM（kill 命令）信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 设定超时时间
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 关闭服务
	e.Shutdown(ctx)
	return nil
}

// Shutdown 关闭服务
func (e *Engine) Shutdown(ctx context.Context) {
	for _, server := range e.servers {
		server.Shutdown(ctx)
	}
	e.flags[StepRunning] = false
}

// 添加服务
func (e *Engine) addServer(servers ...serverx.Server) {
	e.runningCheck()
	e.servers = append(e.servers, servers...)
}

// 新增自定义配置器
func (e *Engine) addConfigurator(configurators ...configx.Configurator) {
	e.runningCheck()
	e.configurators = append(e.configurators, configurators...)
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
	e.tablers[source] = append(e.tablers[source], tablers...)
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
