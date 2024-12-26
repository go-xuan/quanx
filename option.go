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
		e.opts[customPort] = true
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

// SetConfigurator 自定义配置器
func SetConfigurator(configurators ...configx.Configurator) EngineOptionFunc {
	return func(e *Engine) {
		e.AddConfigurator(configurators...)
	}
}

// SetCustomFunc 设置自定义函数
func SetCustomFunc(funcs ...func()) EngineOptionFunc {
	return func(e *Engine) {
		e.AddCustomFunc(funcs...)
	}
}

// SetGinMiddleware 设置gin中间件
func SetGinMiddleware(funcs ...gin.HandlerFunc) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinMiddleware(funcs...)
	}
}

// SetGinRouter 设置gin的路由加载函数
func SetGinRouter(router ...func(*gin.RouterGroup)) EngineOptionFunc {
	return func(e *Engine) {
		e.AddGinRouter(router...)
	}
}

// SetTable 设置 gormx.Tabler 模型
func SetTable(dst ...gormx.Tabler) EngineOptionFunc {
	return func(e *Engine) {
		e.AddTable(dst...)
	}
}

// SetSourceTable 设置某个数据源的 gormx.Table 模型
func SetSourceTable(source string, dst ...gormx.Tabler) EngineOptionFunc {
	return func(e *Engine) {
		e.AddSourceTable(source, dst...)
	}
}

func EnableDebug() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[enableDebug] = true
	}
}

func EnableNacos() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[enableNacos] = true
	}
}

func EnableQueue() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[enableQueue] = true
	}
}

func MultiDatabase() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[multiDatabase] = true
	}
}

func MultiRedis() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[multiRedis] = true
	}
}

func MultiCache() EngineOptionFunc {
	return func(e *Engine) {
		e.opts[multiCache] = true
	}
}
