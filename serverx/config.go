package serverx

import (
	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/osx"
)

// DefaultConfig 默认服务运行配置
func DefaultConfig() *Config {
	return &Config{
		Name: "quanx-server",
		Host: osx.GetLocalIP(),
		Port: map[string]int{
			HTTP: 8888,
		},
	}
}

// Config 服务运行配置
type Config struct {
	Name string         `json:"name" yaml:"name"` // 服务名称
	Host string         `json:"host" yaml:"host"` // host, 为空时默认获取本地IP
	Port map[string]int `json:"port" yaml:"port"` // 服务端口, 键为服务类型, 值为端口号
}

func (c *Config) Readers() []configx.Reader {
	return []configx.Reader{
		nacosx.NewReader("server.yaml"),
		configx.NewFileReader("server.yaml"),
	}
}

func (c *Config) Valid() bool {
	return c.Name != "" && c.Host != "" && len(c.Port) > 0
}

func (c *Config) Execute() error {
	return nil
}

// Cover 覆盖配置，仅合并非空字段
func (c *Config) Cover(cover *Config) {
	if cover == nil {
		return
	}
	if name := cover.Name; name != "" {
		c.Name = name
	}
	if host := cover.Host; host != "" {
		c.Host = host
	}
	// 合并端口配置
	if cover.Port != nil {
		if c.Port == nil {
			c.Port = make(map[string]int)
		}
		for type_, port := range cover.Port {
			c.Port[type_] = port
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
	if len(c.Port) > 0 && c.Port[HTTP] > 0 {
		return c.Port[HTTP]
	}
	return 0
}

// RegisterServer 注册服务
func (c *Config) RegisterServer() error {
	if nacosx.Initialized() {
		if client := nacosx.GetClient().GetNamingClient(); client != nil {
			InitNacosCenter(nacosx.GetClient().GetGroup(), client)
			if err := Register(c); err != nil {
				return errorx.Wrap(err, "register nacos server instance failed")
			}
		}
	}
	return nil
}
