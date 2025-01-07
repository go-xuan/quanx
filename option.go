package quanx

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/gormx"
)

type Option uint

const (
	enableDebug   Option = iota // debug模式
	enableNacos                 // 启用nacos
	multiDatabase               // 开启多数据源
	multiRedis                  // 开启多redis源
	multiCache                  // 开启多缓存源
	enableQueue                 // 使用队列任务启动
	customPort                  // 自定义端口
	running                     // 正在运行中
)

type EngineOptionFunc = func(e *Engine)

func SetPort(port int) EngineOptionFunc {
	return func(e *Engine) {
		_ = os.Setenv("PORT", strconv.Itoa(port))
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
func SetConfig(cof *Config) EngineOptionFunc {
	return func(e *Engine) {
		e.config = cof
	}
}

// AddConfigurator 自定义配置器
func AddConfigurator(configurators ...configx.Configurator) EngineOptionFunc {
	return func(e *Engine) {
		e.AddConfigurator(configurators...)
	}
}

// AddCustomFunc 设置自定义函数
func AddCustomFunc(funcs ...func()) EngineOptionFunc {
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

// AddTable 设置 gormx.Tabler 模型
func AddTable(dst ...gormx.Tabler) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTable(dst...)
	}
}

// AddSourceTable 设置某个数据源的 gormx.Table 模型
func AddSourceTable(source string, dst ...gormx.Tabler) EngineOptionFunc {
	return func(e *Engine) {
		e.AddSourceTable(source, dst...)
	}
}

func EnableDebug() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableDebug] = true
	}
}

func EnableNacos() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableNacos] = true
	}
}

func EnableQueue() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[enableQueue] = true
	}
}

func MultiDatabase() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[multiDatabase] = true
	}
}

func MultiRedis() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[multiRedis] = true
	}
}

func MultiCache() EngineOptionFunc {
	return func(e *Engine) {
		e.switches[multiCache] = true
	}
}