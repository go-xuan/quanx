package serverx

import (
	"github.com/go-xuan/utilx/osx"
)

// DefaultConfig 默认服务配置
func DefaultConfig() *Config {
	return &Config{
		Name: "quanx-server",
		Host: osx.GetLocalIP(),
		Port: make(map[string]int),
	}
}

// Config 服务配置
type Config struct {
	Name string         `json:"name" yaml:"name"` // 服务名称
	Host string         `json:"host" yaml:"host"` // host, 为空时默认获取本地IP
	Port map[string]int `json:"port" yaml:"port"` // 服务端口
}

// Cover 覆盖配置
func (c *Config) Cover(config *Config) {
	if config == nil {
		return
	}
	if config.Name != "" {
		c.Name = config.Name
	}
	if config.Host != "" {
		c.Host = config.Host
	}
	// 合并端口配置
	if config.Port != nil {
		if c.Port == nil {
			c.Port = make(map[string]int)
		}
		for k, v := range config.Port {
			c.Port[k] = v
		}
	}
}

// GetName 获取服务名
func (c *Config) GetName() string {
	return c.Name
}

// GetHost 获取服务host
func (c *Config) GetHost() string {
	if c.Host == "" {
		c.Host = osx.GetLocalIP()
	}
	return c.Host
}

// GetPort 获取服务端口
func (c *Config) GetPort() int {
	if len(c.Port) > 0 && c.Port["http"] > 0 {
		return c.Port["http"]
	}
	return 0
}
