package quanx

import (
	"context"

	"github.com/go-xuan/utilx/taskx"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/serverx"
)

// EngineOption 配置选项
type EngineOption = func(e *Engine)

// InitHttpConfig 初始化http服务配置
func InitHttpConfig(server *serverx.Config) EngineOption {
	return func(e *Engine) {
		e.config.Server = server
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

// AddServer 添加服务
func AddServer(server ...serverx.Server) EngineOption {
	return func(e *Engine) {
		e.addServer(server...)
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
func AddTaskBefore(baseStep, name string, task func(context.Context) error) EngineOption {
	return func(e *Engine) {
		e.addTaskBefore(baseStep, name, task)
	}
}

// AddTaskAfter 后插队添加任务
func AddTaskAfter(baseStep, name string, task func(context.Context) error) EngineOption {
	return func(e *Engine) {
		e.addTaskAfter(baseStep, name, task)
	}
}

// ReadLocalConfig 读取nacos配置（需保证nacos已经提前初始化，所以以自定义函数的形式延迟执行）
func ReadLocalConfig(config any, path string) EngineOption {
	return func(e *Engine) {
		e.readLocalConfig(config, path)
	}
}

// ReadNacosConfig 读取nacos配置（需保证nacos已经提前初始化，所以以自定义函数的形式延迟执行）
func ReadNacosConfig(config any, dataId string, listen ...bool) EngineOption {
	return func(e *Engine) {
		e.readNacosConfig(config, dataId, listen...)
	}
}
