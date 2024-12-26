package quanx

import (
	"strings"

	"github.com/go-xuan/quanx/core/cachex"
	"github.com/go-xuan/quanx/core/gormx"
	"github.com/go-xuan/quanx/core/nacosx"
	"github.com/go-xuan/quanx/core/redisx"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/types/stringx"
)

// GetServer 获取服务配置
func GetServer() *Server {
	return GetEngine().config.Server
}

// Config 服务配置
type Config struct {
	Server   *Server             `yaml:"server"`   // 服务配置
	Log      *logx.Config        `yaml:"log"`      // 日志配置
	Nacos    *nacosx.Config      `yaml:"nacos"`    // nacos访问配置
	Database *gormx.MultiConfig  `yaml:"database"` // 数据源配置
	Redis    *redisx.MultiConfig `yaml:"redis"`    // redis配置
	Cache    *cachex.MultiConfig `yaml:"cache"`    // 缓存配置
}

// Server 服务配置
type Server struct {
	Name   string `yaml:"name" default:"app"`       // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix"`                   // api prefix（接口根路由）
}

// ApiPrefix API路由前缀
func (s *Server) ApiPrefix() string {
	prefix := stringx.IfZero(s.Prefix, s.Name)
	return stringx.AddPrefix(strings.ToLower(prefix), "/")
}

// Instance 服务实例
func (s *Server) Instance() nacosx.ServerInstance {
	return nacosx.ServerInstance{
		Name: s.Name,
		Host: s.Host,
		Port: s.Port,
	}
}
