package engine

import (
	"github.com/go-xuan/quanx/public/gormx"
	"github.com/go-xuan/quanx/public/logx"
	"github.com/go-xuan/quanx/public/nacosx"
	"github.com/go-xuan/quanx/public/redisx"
)

// 应用配置
type Config struct {
	Server   *Server        `yaml:"server"`   // 服务配置
	Log      *logx.Config   `yaml:"log"`      // 日志配置
	Nacos    *nacosx.Config `yaml:"nacos"`    // nacos访问配置
	Database gormx.Configs  `yaml:"database"` // 数据源配置
	Redis    redisx.Configs `yaml:"redis"`    // redis配置
}

// 服务配置
type Server struct {
	Name        string `yaml:"name"`                     // 服务名
	Host        string `yaml:"host" default:"127.0.0.1"` // 服务host
	Port        int    `yaml:"port" default:"8888"`      // 服务端口
	Env         string `yaml:"env" default:"localhost"`  // server run environment
	Prefix      string `yaml:"prefix" default:"api"`     // prefix of api (RESTFul api prefix)
	SessionName string `yaml:"sessionName" default:"my_session"`
}
