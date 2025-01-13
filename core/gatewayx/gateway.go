package gatewayx

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-xuan/quanx/core/configx"
	"github.com/go-xuan/quanx/core/nacosx"
	"github.com/go-xuan/quanx/os/errorx"
	"github.com/go-xuan/quanx/types/stringx"
)

var Gateway = &Servers{}

type Servers []*Server

func (s *Servers) Format() string {
	if s != nil && len(*s) > 0 {
		b, _ := json.Marshal(s)
		return string(b)
	}
	return ""
}

func (s *Servers) Reader() *configx.Reader {
	return &configx.Reader{
		FilePath:    "gateway.yaml",
		NacosDataId: "gateway.yaml",
		Listen:      true,
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
	if err := ListenConfigChanged(group, dataId); err != nil {
		return "", "", errorx.Wrap(err, "监听微服务网关配置失败")
	}
	var auth string
	for _, server := range *Gateway {
		if MatchUrl(url, server.Router) {
			if len(server.Ignore) > 0 {
				for _, item := range server.Ignore {
					if stringx.Index(url, strings.TrimSpace(item)) >= 0 {
						auth = server.Auth
						break
					}
				}
			}
			if addr, err := nacosx.SelectOneHealthyInstance(server.Name, server.Group); err != nil {
				return "", "", errorx.Wrap(err, "微服务实例未注册")
			} else {
				addr = "http://" + addr + server.Prefix
				return addr, auth, nil
			}
		}
	}
	return "", "", errorx.Errorf("未找到对应的网关路由配置，请检查微服务配置文件，或者确认请求接口[%s]是否正确", url)

}

// ListenConfigChanged 校验nacos配置
func ListenConfigChanged(group, dataId string) error {
	if data, ok := nacosx.GetConfigMonitor().Get(group, dataId); ok {
		// 将当前最新的content数据同步到servers
		if err := data.Unmarshal(Gateway); err != nil {
			return errorx.Wrap(err, "unmarshal error")
		}
		// 更新nacos监控中配置值
		data.SetChanged(false)
	}
	return nil
}

// MatchUrl URL匹配
func MatchUrl(uri, rule string) bool {
	if rule == "*" || rule == "/*" {
		return true
	} else if stringx.Index(rule, ",") >= 0 {
		var list = strings.Split(rule, ",")
		for _, item := range list {
			if stringx.Index(uri, strings.TrimSpace(item)) >= 0 {
				return true
			}
		}
	} else if strings.HasSuffix(rule, `/**`) {
		prefix := rule[:len(rule)-3]
		return strings.HasPrefix(uri, prefix)
	} else if strings.HasSuffix(rule, `/*`) {
		prefix := rule[:len(rule)-2]
		if strings.HasPrefix(uri, prefix) {
			uri = uri[len(prefix):]
			return stringx.Index(uri, `/`) < 0
		}
	} else {
		return stringx.Index(uri, rule) >= 0
	}
	return false
}
