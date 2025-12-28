package appx

import (
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/constx"
	"github.com/go-xuan/quanx/serverx"
)

// Option 配置选项
type Option = func(e *Engine)

// SetServerConfig 设置预制服务配置
func SetServerConfig(server *serverx.Config) Option {
	return func(e *Engine) {
		e.config.Server = server
	}
}

// AddConfigurator 添加自定义配置器
func AddConfigurator(configurators ...configx.Configurator) Option {
	return func(e *Engine) {
		e.configurators = append(e.configurators, configurators...)
	}
}

// AddTable 添加表结构（默认数据源）
func AddTable(tablers ...interface{}) Option {
	return AddSourceTable(constx.DefaultSource, tablers...)
}

// AddSourceTable 添加表结构（指定数据源）
func AddSourceTable(source string, tablers ...interface{}) Option {
	return func(e *Engine) {
		e.tablers[source] = append(e.tablers[source], tablers...)
	}
}

// AddServer 添加服务
func AddServer(servers ...serverx.Server) Option {
	return func(e *Engine) {
		e.AddServer(servers...)
	}
}
