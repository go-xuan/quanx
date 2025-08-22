package serverx

import (
	"fmt"
	"strings"

	"github.com/go-xuan/utilx/stringx"
)

// Config 服务配置
type Config struct {
	Name   string `json:"name" yaml:"name" default:"app"`     // 服务名
	IP     string `json:"ip" yaml:"ip"`                       // 服务IP
	Port   int    `json:"port" yaml:"port" default:"8888"`    // 服务端口
	Prefix string `json:"prefix" yaml:"prefix"`               // 服务根路由
	Domain string `json:"domain" yaml:"domain"`               // 服务域名
	Debug  bool   `json:"debug" yaml:"debug" default:"false"` // debug模式
	Group  string `json:"group" yaml:"group"`                 // 服务分组
}

// ApiPrefix API路由前缀
func (c *Config) ApiPrefix() string {
	prefix := stringx.IfZero(c.Prefix, c.Name) // 优先使用配置的前缀, 否则使用服务名
	prefix = stringx.ToSnake(prefix)           // 转为蛇形命名
	return stringx.AddPrefix(strings.ToLower(prefix), "/")
}

// ApiDomain API域名
func (c *Config) ApiDomain() string {
	if c.Domain != "" {
		return c.Domain
	}
	return fmt.Sprintf("http://%s:%d", c.IP, c.Port)
}

// NacosGroup Nacos配置分组
func (c *Config) NacosGroup() string {
	if c.Group != "" {
		return c.Group
	}
	return c.Name
}
