package quanx

import (
	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/extra/configx"
)

type Option uint

const (
	enableDebug Option = iota // debug模式
	customPort                // 自定义端口
	running                   // 正在运行中
)

type EngineOptionFunc = func(e *Engine)

func SetPort(port int) EngineOptionFunc {
	return func(e *Engine) {
		if server := e.config.Server; server == nil {
			e.config.Server = &Server{Port: port}
		} else {
			server.Port = port
		}
		e.switches[customPort] = true
	}
}

// SetConfigDir 设置配置文件
func SetConfigDir(dir string) EngineOptionFunc {
	return func(e *Engine) {
		e.SetConfigDir(dir)
	}
}

// SetConfig 自定义应用配置
func SetConfig(config *Config) EngineOptionFunc {
	return func(e *Engine) {
		e.config = config
	}
}

// AddConfigurator 自定义配置器
func AddConfigurator(configurators ...configx.Configurator) EngineOptionFunc {
	return func(e *Engine) {
		e.AddConfigurator(configurators...)
	}
}

// AddCustomFunc 设置自定义函数
func AddCustomFunc(funcs ...func() error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddCustomFunc(funcs...)
	}
}

// AddGinMiddleware 设置gin中间件
func AddGinMiddleware(funcs ...gin.HandlerFunc) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinMiddleware(funcs...)
	}
}

// AddGinRouter 设置gin的路由加载函数
func AddGinRouter(router ...func(*gin.RouterGroup)) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinRouter(router...)
	}
}

// AddTable 添加表结构
func AddTable(tablers ...interface{}) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTable(tablers...)
	}
}

// AddSourceTable 添加数据源表结构
func AddSourceTable(source string, tablers ...interface{}) EngineOptionFunc {
	return func(e *Engine) {
		e.AddSourceTable(source, tablers...)
	}
}

// AddQueueTask 使用后，会自动以队列方式来启动服务
func AddQueueTask(name, before string, task func() error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTaskBefore(name, before, task)
	}
}

// EnableDebug 启用debug
func EnableDebug() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableDebug] = true
	}
}
