package gatewayx

import (
	"github.com/go-xuan/utilx/stringx"
)

// Proxy 服务代理
type Proxy struct {
	Url    string `yaml:"url" json:"url"`       // URL规则
	Server string `yaml:"server" json:"server"` // 代理服务
	Prefix string `yaml:"prefix" json:"prefix"` // 代理服务url前缀
}

// Match URL匹配服务
func (p *Proxy) Match(url string) bool {
	if stringx.MatchUrl(url, p.Url) {
		return true
	}
	return false
}
