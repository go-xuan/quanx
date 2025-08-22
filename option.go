package quanx

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/serverx"
)

type Option uint

const (
	enableDebug Option = iota // debug模式
	customPort                // 自定义端口
	running                   // 正在运行中
)

// EngineOptionFunc 配置选项
type EngineOptionFunc = func(e *Engine)

func SetPort(port int) EngineOptionFunc {
	return func(e *Engine) {
		if server := e.config.Server; server == nil {
			e.config.Server = &serverx.Config{Port: port}
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
func AddTaskBefore(base, name string, task func(context.Context) error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTaskBefore(base, name, task)
	}
}

// AddTaskAfter 后插队添加任务
func AddTaskAfter(base, name string, task func(context.Context) error) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTaskAfter(base, name, task)
	}
}

// Debug 启用debug
func Debug() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableDebug] = true
	}
}
