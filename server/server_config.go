package server

import (
	"github.com/quanxiaoxuan/quanx/middleware/gormx"
	"github.com/quanxiaoxuan/quanx/middleware/hugegraphx"
	"github.com/quanxiaoxuan/quanx/middleware/logx"
	"github.com/quanxiaoxuan/quanx/middleware/nacosx"
	"github.com/quanxiaoxuan/quanx/middleware/redisx"
)

// 应用配置
type Config struct {
	Server    Server                `yaml:"server"`    // 服务配置
	Log       logx.Config           `yaml:"log"`       // 日志配置
	Nacos     nacosx.Config         `yaml:"nacos"`     // nacos访问配置
	Configs   nacosx.ConfigItemList `yaml:"configs"`   // nacos配置清单
	Database  gormx.Config          `yaml:"database"`  // 数据库访问配置
	Redis     redisx.Config         `yaml:"redis"`     // redis配置
	Hugegraph hugegraphx.Config     `yaml:"hugegraph"` // hugegraph配置
}

// 服务配置
type Server struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host" default:"127.0.0.1"`
	Port   string `yaml:"port" default:"8888"`
	Env    string `yaml:"env" default:"localhost"` // server run environment
	Prefix string `yaml:"prefix" default:"api"`    // prefix of api (RESTFul api prefix)
}
