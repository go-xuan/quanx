package serverx

import (
	"fmt"
	"strings"

	"github.com/go-xuan/utilx/osx"
	"github.com/go-xuan/utilx/stringx"
)

// Config 服务配置
type Config struct {
	Name   string `json:"name" yaml:"name" default:"app"`     // 服务名
	Group  string `json:"group" yaml:"group"`                 // 服务分组
	IP     string `json:"ip" yaml:"ip"`                       // 服务IP
	Port   int    `json:"port" yaml:"port" default:"8888"`    // 服务端口
	Prefix string `json:"prefix" yaml:"prefix"`               // 服务根路由
	Domain string `json:"domain" yaml:"domain"`               // 服务域名
	Debug  bool   `json:"debug" yaml:"debug" default:"false"` // debug模式
}

// ApiPrefix API路由前缀
func (c *Config) ApiPrefix() string {
	prefix := stringx.IfZero(c.Prefix, c.Name) // 优先使用配置的前缀, 否则使用服务名
	prefix = stringx.ToSnake(prefix)           // 转为蛇形命名
	return stringx.AddPrefix(strings.ToLower(prefix), "/")
}

// GetIP 获取服务IP
func (c *Config) GetIP() string {
	if c.IP == "" {
		c.IP = osx.GetLocalIP()
	}
	return c.IP
}

// GetAddress 获取服务域名
func (c *Config) GetAddress() string {
	if c.Domain == "" {
		c.Domain = fmt.Sprintf("http://%s:%d", c.GetIP(), c.Port)
	}
	return c.Domain
}

// GetGroup 获取分组
func (c *Config) GetGroup() string {
	if c.Group == "" {
		c.Group = c.Name
	}
	return c.Group
}
