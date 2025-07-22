package gatewayx

import (
	"strings"

	"github.com/go-xuan/utilx/stringx"
)

// Proxy 服务代理
type Proxy struct {
	Url    string   `yaml:"url" json:"url"`       // URL规则
	Server string   `yaml:"server" json:"server"` // 代理服务
	Prefix string   `yaml:"prefix" json:"prefix"` // 代理服务url前缀
	Auth   string   `yaml:"auth" json:"auth"`     // 代理服务鉴权方式(cookie/token/no)
	Ignore []string `yaml:"ignore" json:"ignore"` // 代理服务鉴权白名单(白名单)
}

// Match URL匹配服务
func (p *Proxy) Match(url string) bool {
	if stringx.MatchUrl(url, p.Url) {
		return true
	}
	return false
}

// AuthCheck 鉴权检查
func (p *Proxy) AuthCheck(url string) string {
	if p.Auth != "" && len(p.Ignore) > 0 {
		for _, ignore := range p.Ignore {
			if stringx.Index(url, strings.TrimSpace(ignore)) >= 0 {
				return ""
			}
		}
	}
	return p.Auth
}
