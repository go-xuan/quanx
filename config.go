package quanx

import (
	"github.com/go-xuan/quanx/cachex"
	"github.com/go-xuan/quanx/gormx"
	"github.com/go-xuan/quanx/logx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/quanx/redisx"
	"github.com/go-xuan/quanx/serverx"
)

// Config 服务配置
type Config struct {
	Server   *serverx.Config     `json:"server" yaml:"server"`     // 服务配置
	Log      *logx.Config        `json:"log" yaml:"log"`           // 日志配置
	Nacos    *nacosx.Config      `json:"nacos" yaml:"nacos"`       // nacos访问配置
	Database *gormx.MultiConfig  `json:"database" yaml:"database"` // 数据源配置
	Redis    *redisx.MultiConfig `json:"redis" yaml:"redis"`       // redis配置
	Cache    *cachex.MultiConfig `json:"cache" yaml:"cache"`       // 缓存配置
}
