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

// AddCustomFunc 添加自定义函数
func AddCustomFunc(funcs ...func() error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddCustomFunc(funcs...)
	}
}

// AddGinMiddleware 添加gin中间件
func AddGinMiddleware(middleware ...gin.HandlerFunc) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinMiddleware(middleware...)
	}
}

// AddGinRouter 添加gin的路由加载函数
func AddGinRouter(router ...func(*gin.RouterGroup)) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinRouter(router...)
	}
}

// AddTable 添加表结构（默认数据源）
func AddTable(tablers ...interface{}) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTable(tablers...)
	}
}

// AddSourceTable 添加表结构（指定数据源）
func AddSourceTable(source string, tablers ...interface{}) EngineOptionFunc {
	return func(e *Engine) {
		e.AddSourceTable(source, tablers...)
	}
}

// AddTaskBefore 前插队添加任务
func AddTaskBefore(name, before string, task func() error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTaskBefore(name, before, task)
	}
}

// AddTaskAfter 后插队添加任务
func AddTaskAfter(name, after string, task func() error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTaskAfter(name, after, task)
	}
}

// EnableDebug 启用debug
func EnableDebug() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableDebug] = true
	}
}
