package app

import (
	"fmt"
	"github.com/go-xuan/quanx/app/nacosx"
	"github.com/go-xuan/quanx/os/cachex"
	"github.com/go-xuan/quanx/os/logx"
	"github.com/go-xuan/quanx/server/gormx"
	"github.com/go-xuan/quanx/server/redisx"
	"github.com/go-xuan/quanx/types/stringx"
)

// GetServer 获取服务配置
func GetServer() *Server {
	return GetEngine().config.Server
}

// Config 服务配置
type Config struct {
	Server   *Server              `yaml:"server"`   // 服务配置
	Log      *logx.Log            `yaml:"log"`      // 日志配置
	Nacos    *nacosx.Nacos        `yaml:"nacos"`    // nacos访问配置
	Database *gormx.MultiDatabase `yaml:"database"` // 数据源配置
	Redis    *redisx.MultiRedis   `yaml:"redis"`    // redis配置
	Cache    *cachex.MultiCache   `yaml:"cache"`    // 缓存配置
}

// Server 服务配置
type Server struct {
	Name   string `yaml:"name"`                     // 服务名
	Host   string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port   int    `yaml:"port" default:"8888"`      // 服务端口
	Prefix string `yaml:"prefix"`                   // api prefix（接口根路由）
	Debug  bool   `yaml:"debug" default:"false"`    // 是否调试环境
}

// ApiPrefix API路由前缀
func (s *Server) ApiPrefix() string {
	prefix := stringx.IfZero(s.Prefix, s.Name)
	return stringx.AddPrefix(prefix, "/")
}

// ApiHost 服务host
func (s *Server) ApiHost() string {
	return fmt.Sprintf(`http://%s:%d`, s.Host, s.Port)
}

// Instance 服务实例
func (s *Server) Instance() nacosx.ServerInstance {
	return nacosx.ServerInstance{Name: s.Name, Host: s.Host, Port: s.Port}
}
