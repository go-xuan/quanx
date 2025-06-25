package gatewayx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-xuan/utilx/errorx"
	"github.com/go-xuan/utilx/stringx"

	"github.com/go-xuan/quanx/configx"
	"github.com/go-xuan/quanx/nacosx"
)

var Gateway = &Servers{}

type Servers []*Server

func (s *Servers) Info() string {
	if s != nil && len(*s) > 0 {
		b, _ := json.Marshal(s)
		return string(b)
	}
	return ""
}

func (*Servers) Reader(from configx.From) configx.Reader {
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

func (s *Servers) Execute() error {
	return nil
}

// Server 微服务网关配置
type Server struct {
	Name   string   `yaml:"name" json:"name"`     // 微服务名称
	Group  string   `yaml:"group" json:"group"`   // 微服务分组
	Prefix string   `yaml:"prefix" json:"prefix"` // 微服务API前缀
	Router string   `yaml:"router" json:"router"` // API路由
	Auth   string   `yaml:"auth" json:"auth"`     // 鉴权方式(cookie/token/no)
	Ignore []string `yaml:"ignore" json:"ignore"` // 忽略鉴权(白名单)
}

func (s Server) Format() string {
	return fmt.Sprintf("name=%s group=%s prefix=%s router=%s auth=%s",
		s.Name, s.Group, s.Prefix, s.Router, s.Auth)
}

// GetServerProxyAddr 获取微服务addr
func GetServerProxyAddr(group, dataId, url string) (string, string, error) {
	var auth string
	for _, server := range *Gateway {
		if stringx.MatchUrl(url, server.Router) {
			if len(server.Ignore) > 0 {
				for _, item := range server.Ignore {
					if stringx.Index(url, strings.TrimSpace(item)) >= 0 {
						auth = server.Auth
						break
					}
				}
			}
			if instance, err := nacosx.SelectOneHealthyInstance(server.Name, server.Group); err != nil {
				return "", "", errorx.Wrap(err, "微服务实例未注册")
			} else {
				addr := fmt.Sprintf("http://%s:%d/%s",
					instance.Host, instance.Port, strings.TrimPrefix(server.Prefix, "/"))
				return addr, auth, nil
			}
		}
	}
	return "", "", errorx.Errorf("未找到对应的网关路由配置，请检查微服务配置文件，或者确认请求接口[%s]是否正确", url)

}
