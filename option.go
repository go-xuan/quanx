package quanx

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/utilx/taskx"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/serverx"
)

// EngineOption 配置选项
type EngineOption = func(e *Engine)

// SetPort 设置服务端口
func SetPort(port int) EngineOption {
	return func(e *Engine) {
		if server := e.config.Server; server == nil {
			e.config.Server = &serverx.Config{
				Port: port,
			}
		} else {
			server.Port = port
		}
	}
}

// Debug 启用debug
func Debug() EngineOption {
	return func(e *Engine) {
		if server := e.config.Server; server == nil {
			e.config.Server = &serverx.Config{
				Debug: true,
			}
		} else {
			server.Debug = true
		}
	}
}

// SetConfig 自定义应用配置
func SetConfig(config *Config) EngineOption {
	return func(e *Engine) {
		e.config = config
	}
}

// AddConfigurator 自定义配置器
func AddConfigurator(configurators ...configx.Configurator) EngineOption {
	return func(e *Engine) {
		e.addConfigurator(configurators...)
	}
}

// AddExecute 添加自定义函数
func AddExecute(executes ...taskx.Execute) EngineOption {
	return func(e *Engine) {
		e.addExecute(executes...)
	}
}

// AddGinRouter 添加gin的路由加载函数
func AddGinRouter(router ...func(*gin.RouterGroup)) EngineOption {
	return func(e *Engine) {
		e.addGinRouter(router...)
	}
}

// AddGinMiddleware 添加gin中间件
func AddGinMiddleware(middleware ...gin.HandlerFunc) EngineOption {
	return func(e *Engine) {
		e.addGinMiddleware(middleware...)
	}
}

// AddTable 添加表结构（默认数据源）
func AddTable(tablers ...interface{}) EngineOption {
	return func(e *Engine) {
		e.addTable(tablers...)
	}
}

// AddSourceTable 添加表结构（指定数据源）
func AddSourceTable(source string, tablers ...interface{}) EngineOption {
	return func(e *Engine) {
		e.addSourceTable(source, tablers...)
	}
}

// AddTaskBefore 前插队添加任务
func AddTaskBefore(base, name string, task func(context.Context) error) EngineOption {
	return func(e *Engine) {
		e.addTaskBefore(base, name, task)
	}
}

// AddTaskAfter 后插队添加任务
func AddTaskAfter(base, name string, task func(context.Context) error) EngineOption {
	return func(e *Engine) {
		e.addTaskAfter(base, name, task)
	}
}

// ReadLocalConfig 读取nacos配置（需保证nacos已经提前初始化，所以以自定义函数的形式延迟执行）
func ReadLocalConfig(config any, path string) EngineOption {
	return func(e *Engine) {
		e.readLocalConfig(config, path)
	}
}

// ReadNacosConfig 读取nacos配置（需保证nacos已经提前初始化，所以以自定义函数的形式延迟执行）
func ReadNacosConfig(v any, dataId string, listen ...bool) EngineOption {
	return func(e *Engine) {
		e.readNacosConfig(v, dataId, listen...)
	}
}
