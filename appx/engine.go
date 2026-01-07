package appx

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/dbx"
	"github.com/go-xuan/quanx/serverx"
)

const (
	DefaultConfigName = "config.yaml"

	FlagInit    = "init"    // 初始化标识
	FlagRunning = "running" // 服务运行标识
)

var (
	engine *Engine   // 服务启动器
	once   sync.Once // 单例模式
)

// NewEngine 初始化Engine
func NewEngine(options ...Option) *Engine {
	once.Do(func() {
		engine = &Engine{
			config:        &Config{},
			configurators: make([]configx.Configurator, 0),
			tablers:       make(map[string][]interface{}),
			servers:       make([]serverx.Server, 0),
			flags:         make(map[string]bool),
		}
	})
	// 添加配置选项
	for _, option := range options {
		option(engine)
	}
	return engine
}

// GetEngine 获取当前Engine
func GetEngine() *Engine {
	if engine == nil {
		engine = NewEngine(
			SetServerConfig(serverx.DefaultConfig()),
		)
	}
	return engine
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	return GetEngine().config
}

// Engine 应用启动Engine
type Engine struct {
	config        *Config                  // 服务启动配置
	configurators []configx.Configurator   // 配置器
	tablers       map[string][]interface{} // 初始化表结构
	servers       []serverx.Server         // http/grpc或者其他服务
	flags         map[string]bool          // 服务运行标识
}

// AddServer 添加服务
func (e *Engine) AddServer(servers ...serverx.Server) {
	e.servers = append(e.servers, servers...)
}

// MustInit 应用初始化，确保必须且仅初始化一次
func (e *Engine) MustInit(ctx context.Context) {
	errorx.Panic(e.initOnce(ctx))
}

// LoadConfigurator 初始化自定义配置器
func (e *Engine) LoadConfigurator(configurators ...configx.Configurator) {
	for _, configurator := range configurators {
		errorx.Panic(configx.LoadConfigurator(configurator))
	}
}

// InitTable 初始化数据库表结构以及数据
func (e *Engine) InitTable(source string, tablers ...interface{}) {
	client := dbx.GetClient(source)
	errorx.Panic(dbx.InitGormTable(client, tablers...))
}

// Shutdown 关闭服务
func (e *Engine) Shutdown(ctx context.Context) {
	serverx.Shutdown(ctx, e.servers...)
	e.reset()
	log.WithContext(ctx).Info("shutdown complete")
}

// RUN 运行应用
func (e *Engine) RUN(ctx context.Context) {
	e.checkRunning()   // 检查服务是否已运行
	e.MustInit(ctx)    // 初始化应用
	e.startServer(ctx) // 启动服务
	e.keepRunning(ctx) // 保持服务运行
	e.Shutdown(ctx)    // 关闭服务
}

// 检查服务运行状态
func (e *Engine) checkRunning() {
	if e.flags[FlagRunning] {
		panic("engine has already running")
	}
}

// 应用初始化，确保必须且仅初始化一次
func (e *Engine) initOnce(_ context.Context) error {
	if e.flags[FlagInit] {
		return nil
	}
	// 初始化应用配置（日志、nacos、数据库、redis、缓存等）
	reader := configx.NewFileReader(DefaultConfigName)
	if err := e.config.Init(reader); err != nil {
		return errorx.Wrap(err, "init default config error")
	}
	// 初始化配置器
	for _, configurator := range e.configurators {
		if err := configx.LoadConfigurator(configurator); err != nil {
			return errorx.Wrap(err, "load configurators error")
		}
	}
	// 初始化数据库表结构
	if err := dbx.InitTable(e.tablers); err != nil {
		return errorx.Wrap(err, "init table error")
	}
	// 设置初始化标识
	e.flags[FlagInit] = true
	return nil
}

// 启动服务
func (e *Engine) startServer(ctx context.Context) {
	if config := e.config.Server; config != nil {
		errorx.Panic(serverx.Start(ctx, config, e.servers...))
	}
}

// 保持服务运行，等待信号量关闭服务
func (e *Engine) keepRunning(ctx context.Context) {
	e.flags[FlagRunning] = true

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	// 设定超时时间，确保服务有足够时间关闭
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
}

// 重置应用
func (e *Engine) reset() {
	e.configurators = make([]configx.Configurator, 0)
	e.tablers = make(map[string][]interface{})
	e.servers = make([]serverx.Server, 0)
	e.flags = make(map[string]bool)
}
