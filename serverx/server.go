package serverx

import (
	"context"

	"github.com/go-xuan/utilx/osx"
)

// Server 服务
type Server interface {
	Run(config *Config) error
	Shutdown(context.Context)
	IsRunning() bool
}

// DefaultConfig 默认服务配置
func DefaultConfig() *Config {
	return &Config{
		Name:  "quanx-server",
		Debug: false,
	}
}

// Config 服务配置
type Config struct {
	Name  string `json:"name" yaml:"name"`   // 名称
	Host  string `json:"host" yaml:"host"`   // host
	Debug bool   `json:"debug" yaml:"debug"` // 开启调试模式
	Http  int    `json:"http" yaml:"http"`   // http服务端口, 默认8080
	Grpc  int    `json:"grpc" yaml:"grpc"`   // grpc服务端口, 0表示不开启grpc服务
}

// GetName 获取服务名
func (a *Config) GetName() string {
	return a.Name
}

// GetHost 获取服务host
func (a *Config) GetHost() string {
	if a.Host == "" {
		a.Host = osx.GetLocalIP()
	}
	return a.Host
}

// GetPort 获取服务端口
func (a *Config) GetPort() int {
	if a.Http != 0 {
		return a.Http
	}
	return 9999
}
