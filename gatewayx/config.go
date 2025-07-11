package gatewayx

import (
	"github.com/go-xuan/utilx/errorx"
	log "github.com/sirupsen/logrus"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

var Gateway = &Config{}

// Config 网关配置
type Config struct {
	Proxies []*Proxy `yaml:"proxies" json:"proxies"` // 服务代理
}

func (c *Config) NeedRead() bool {
	return len(c.Proxies) == 0
}

func (*Config) Reader(from configx.From) configx.Reader {
	switch from {
	case configx.FromNacos:
		return &nacosx.Reader{
			DataId: "gateway.yaml",
			Listen: true,
		}
	case configx.FromFile:
		return &configx.FileReader{
			Name: "gateway.yaml",
		}
	default:
		return nil
	}
}

func (c *Config) Execute() error {
	if len(c.Proxies) == 0 {
		return errorx.New("gateway not initialized, gateway.yaml is invalid")
	}
	var fields = log.Fields{}
	for _, proxy := range c.Proxies {
		fields[proxy.Url] = proxy.Server
	}
	log.WithFields(fields).Info("gateway init success")
	return nil
}

func (c *Config) GetProxy(url string) *Proxy {
	for _, proxy := range c.Proxies {
		if proxy.Match(url) {
			return proxy
		}
	}
	return nil
}
